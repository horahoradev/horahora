package grpcserver

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/horahoradev/horahora/video_service/storage"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/horahoradev/horahora/video_service/internal/dashutils"

	"github.com/horahoradev/horahora/video_service/internal/models"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	proto "github.com/horahoradev/horahora/video_service/protocol"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	_ "github.com/aws/aws-sdk-go-v2/aws/external"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/google/uuid"
)

const uploadDir = "/tmp/"

var _ proto.VideoServiceServer = (*GRPCServer)(nil)

type GRPCServer struct {
	VideoModel *models.VideoModel
	Local      bool
	OriginFQDN string
	Storage    storage.Storage
	proto.UnsafeVideoServiceServer
}

// TODO: API is getting bloated
func NewGRPCServer(bucketName string, db *sqlx.DB, port int, originFQDN string, local bool,
	redisClient *redis.Client, client userproto.UserServiceClient, tracer opentracing.Tracer,
	storageBackend, apiID, apiKey string, approvalThreshold int, minioEndpoint string) error {
	g, err := initGRPCServer(bucketName, db, client, local, redisClient, originFQDN, storageBackend, apiID, apiKey, approvalThreshold, minioEndpoint)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TODO: context and return
	go g.transcodeAndUploadVideos()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(tracer))))
	proto.RegisterVideoServiceServer(grpcServer, g)
	return grpcServer.Serve(lis)
}

func initGRPCServer(bucketName string, db *sqlx.DB, client userproto.UserServiceClient, local bool,
	redisClient *redis.Client, originFQDN, storageBackend, apiID, apiKey string, approvalThreshold int, minioEndpoint string) (*GRPCServer, error) {

	g := &GRPCServer{
		Local:      local,
		OriginFQDN: originFQDN,
	}

	var err error

	switch storageBackend {
	case "b2":
		g.Storage, err = storage.NewB2(apiID, apiKey, bucketName)
		if err != nil {
			return nil, err
		}
	case "s3":
		g.Storage, err = storage.NewS3(bucketName)
	case "minio":
		g.Storage, err = storage.NewMinio(minioEndpoint, apiID, apiKey, bucketName)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown storage backend %s", storageBackend)
	}

	g.VideoModel, err = models.NewVideoModel(db, client, redisClient, approvalThreshold)
	if err != nil {
		return nil, err
	}

	return g, nil
}

type VideoUpload struct {
	Meta         *proto.InputVideoChunk_Meta
	FileData     *os.File
	MetaFileData *os.File
}

func (g GRPCServer) UploadVideo(inpStream proto.VideoService_UploadVideoServer) error {
	log.Info("Handling video upload")

	var video VideoUpload

	// UUID for tmp filename and all uploads provides probabilistic guarantee of uniqueness
	id, err := uuid.NewUUID()
	if err != nil {
		log.Error("Could not generate uuid. Err: %s", err)
		return err
	}

	tmpFile, err := os.Create(uploadDir + id.String())
	if err != nil {
		err = fmt.Errorf("could not create tmp file. Err: %s", err)
		log.Error(err)
		return err
	}

	// This is awkward but it could be worse
	metaTmp, err := os.Create(uploadDir + id.String() + ".json")
	if err != nil {
		err = fmt.Errorf("could not create meta tmp file. Err: %s", err)
		log.Error(err)
		return err
	}

	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())

		metaTmp.Close()
		os.Remove(metaTmp.Name())
	}()

	video.FileData = tmpFile
	video.MetaFileData = metaTmp
loop:
	for {
		chunk, err := inpStream.Recv()
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			err = fmt.Errorf("could not recv. Err: %s", err)
			log.Error(err)
			return err
		}

		switch r := chunk.Payload.(type) {
		case *proto.InputVideoChunk_Content:
			_, err := video.FileData.Write(r.Content.Data)
			if err != nil {
				err = fmt.Errorf("could not write video data to file, err: %s", err)
				log.Error(err)
				return err
			}

		case *proto.InputVideoChunk_Rawmeta:
			// lol
			_, err := video.MetaFileData.Write(r.Rawmeta.Data)
			if err != nil {
				err = fmt.Errorf("could not write metadata data to file, err: %s", err)
				log.Error(err)
				return err
			}

		case *proto.InputVideoChunk_Meta:
			if video.Meta != nil {
				err = fmt.Errorf("duplicate metadata in stream")
				log.Error(err)
				return err
			}

			video.Meta = r
			log.Infof("Received metadata for video %s", video.Meta.Meta.Title)
		}
	}

	err = ioutil.WriteFile(video.FileData.Name()+".thumb", video.Meta.Meta.Thumbnail, 0644)
	if err != nil {
		return LogAndRetErr("could not write thumbnail. Err: %s", err)
	}

	log.Infof("Finished receiving file data for %s", video.Meta.Meta.Title)

	// If not local, upload the thumbnail and original video before returning
	if !g.Local {
		// FIXME did it again...
		log.Infof("Uploading thumbnail: %s", video.FileData.Name()+".thumb")
		err = g.Storage.Upload(video.FileData.Name()+".thumb", filepath.Base(video.FileData.Name()+".thumb"))
		if err != nil {
			return err
		}

		// Upload the raw metadata
		log.Infof("Uploading metadata: %s", video.MetaFileData.Name())
		err = g.Storage.Upload(video.MetaFileData.Name(), filepath.Base(video.MetaFileData.Name()))
		if err != nil {
			return err
		}

		// Upload the original video
		log.Infof("Uploading video: %s", video.FileData.Name())
		err = g.Storage.Upload(video.FileData.Name(), filepath.Base(video.FileData.Name()))
		if err != nil {
			return err
		}
	}

	// This is MESSY
	// thumbnail and original video locations are inferred from the mpd location (which is dumb), so it's written even though
	// the video hasn't been transcoded/chunked and the mpd hasn't been uploaded yet
	// a better solution will be provided in the future... I will fix this... (I'm keeping it backwards compatible for now)
	// TODO: switch to struct for args
	manifestLoc := fmt.Sprintf("%s/%s", g.OriginFQDN, filepath.Base(video.FileData.Name()+".mpd"))

	videoID, err := g.VideoModel.SaveForeignVideo(context.TODO(), video.Meta.Meta.Title, video.Meta.Meta.Description,
		video.Meta.Meta.AuthorUsername, video.Meta.Meta.AuthorUID, video.Meta.Meta.OriginalSite,
		video.Meta.Meta.OriginalVideoLink, video.Meta.Meta.OriginalID, manifestLoc, video.Meta.Meta.Tags, video.Meta.Meta.DomesticAuthorID)
	if err != nil {
		return LogAndRetErr("failed to save video to postgres. Err: %s", err)
	}

	uploadResp := proto.UploadResponse{
		VideoID: videoID,
	}

	log.Infof("Finished handling video %s", video.Meta.Meta.Title)
	return inpStream.SendAndClose(&uploadResp)
}

func LogAndRetErr(fmtStr string, err error) error {
	errWithMsg := fmt.Errorf(fmtStr, err)
	log.Error(errWithMsg)
	return errWithMsg
}

func (g GRPCServer) ForeignVideoExists(ctx context.Context, foreignVideoCheck *proto.ForeignVideoCheck) (*proto.VideoExistenceResponse, error) {
	exists, err := g.VideoModel.ForeignVideoExists(foreignVideoCheck.ForeignVideoID, foreignVideoCheck.ForeignWebsite)
	if err != nil {
		return nil, err
	}

	resp := proto.VideoExistenceResponse{Exists: exists}

	return &resp, nil
}

const NUM_TRANSCODING_WORKERS = 1

// TODO: graceful shutdown or something, lock video acquisition
func (g GRPCServer) transcodeAndUploadVideos() {
	for {
		time.Sleep(time.Second * 10)
		videos, err := g.VideoModel.GetUnencodedVideos()
		if err != nil {
			log.Errorf("could not fetch unencoded videos. Err: %s", err)
			continue
		}

		for _, v := range videos {
			time.Sleep(time.Second * 10)

			func(video models.UnencodedVideo) {
				// TODO: distributed lock goes here

				log.Infof("Transcoding/chunking video id %d uuid %s", video.ID, video.GetMPDUUID())

				vid, err := g.Storage.Fetch(video.GetMPDUUID())

				if err != nil {
					log.Errorf("could not fetch unencoded video id %d from backend. Err: %s", video.ID, err)
					return
				} else {
					defer func() {
						vid.Close()
						os.Remove(vid.Name())
					}()
				}

				s, err := vid.Stat()
				if err != nil {
					log.Errorf("Could not stat video to encode. Err: %s", err)
					return
				}

				if s.Size() >= 1024*1024*200 {
					log.Errorf("Video %d greater than 200mb, skipping", v.ID)
					return
				}

				_, err = vid.Seek(0, 0)
				if err != nil {
					log.Errorf("Could not seek to 0 for video to encode. Err: %s", err)
					return
				}

				transcodeResults, err := dashutils.TranscodeAndGenerateManifest(vid.Name(), g.Local)
				if err != nil {
					err := fmt.Errorf("failed to transcode and chunk. Err: %s", err)
					log.Error(err)
					return
				}

				err = g.UploadMPDSet(transcodeResults)
				if err != nil {
					log.Errorf("failed to upload mpd set. Err: %s", err)
					return
				}

				err = g.VideoModel.MarkVideoAsEncoded(video)
				if err != nil {
					log.Errorf("failed to mark video as encoded. Err: %s", err)
					return
				}

				log.Infof("Video %d has been successfully encoded", video.ID)
			}(v)
		}
	}
}

// UploadMPDSet uploads the files to S3. Files may be overwritten (but they're versioned so they're safe).
// Need to ensure as a precondition that the video hasn't been uploaded before and the temp file ID hasn't been
// used.
func (g GRPCServer) UploadMPDSet(d *dashutils.DASHVideo) error {
	// send manifest to origin
	err := g.Storage.Upload(d.ManifestPath, filepath.Base(d.ManifestPath))
	if err != nil {
		return err
	}

	// Send all of the chunked files
	for _, path := range d.QualityMap {
		err = g.Storage.Upload(path, filepath.Base(path))
		if err != nil {
			return err
		}

		err = os.Remove(path)
		if err != nil {
			log.Error(err)
		}

	}

	return nil
}

// Do we need this?
func (g GRPCServer) DownloadVideo(req *proto.VideoRequest, outputStream proto.VideoService_DownloadVideoServer) error {
	return nil
	// TODO
}

func (g GRPCServer) GetVideoList(ctx context.Context, queryConfig *proto.VideoQueryConfig) (*proto.VideoList, error) {
	switch queryConfig.OrderBy {
	case proto.OrderCategory_rating, proto.OrderCategory_views, proto.OrderCategory_upload_date:
		videos, err := g.VideoModel.GetVideoList(queryConfig.Direction, queryConfig.PageNumber,
			queryConfig.FromUserID, queryConfig.SearchVal, queryConfig.ShowUnapproved, queryConfig.OrderBy)
		if err != nil {
			log.Errorf("Could not get video list. Err: %s", err)
			return nil, err
		}

		numberOfVideos, err := g.VideoModel.GetNumberOfSearchResultsForQuery(queryConfig.FromUserID, queryConfig.SearchVal)
		if err != nil {
			log.Errorf("Could not get count of entries for query. Err: %s", err)
			return nil, err
		}

		return &proto.VideoList{
			Videos:         videos,
			NumberOfVideos: numberOfVideos,
		}, nil

	default:
		st := status.New(codes.InvalidArgument, "invalid order category")
		return nil, st.Err()
	}
}

func (g GRPCServer) RateVideo(ctx context.Context, rating *proto.VideoRating) (*proto.Nothing, error) {
	err := g.VideoModel.AddRatingToVideoID(rating.UserID, rating.VideoID, float64(rating.Rating))
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (g GRPCServer) ViewVideo(ctx context.Context, videoInp *proto.VideoViewing) (*proto.Nothing, error) {
	err := g.VideoModel.IncrementViewsForVideo(videoInp.VideoID)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (g GRPCServer) GetVideo(ctx context.Context, req *proto.VideoRequest) (*proto.VideoMetadata, error) {
	videoMetadata, err := g.VideoModel.GetVideoInfo(req.VideoID)
	if err != nil {
		return nil, err
	}

	return videoMetadata, nil
}

func (g GRPCServer) ApproveVideo(ctx context.Context, req *proto.VideoApproval) (*proto.Nothing, error) {
	if err := g.VideoModel.ApproveVideo(int(req.UserID), int(req.VideoID)); err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (g GRPCServer) MakeComment(ctx context.Context, commentReq *proto.VideoComment) (*proto.Nothing, error) {
	return &proto.Nothing{},
		g.VideoModel.MakeComment(commentReq.UserId, commentReq.VideoId,
			commentReq.ParentComment, commentReq.Comment)

}

func (g GRPCServer) MakeCommentUpvote(ctx context.Context, upvoteReq *proto.CommentUpvote) (*proto.Nothing, error) {
	return &proto.Nothing{}, g.VideoModel.MakeUpvote(upvoteReq.UserId, upvoteReq.CommentId,
		upvoteReq.IsUpvote)
}

func (g GRPCServer) GetCommentsForVideo(ctx context.Context, commentListReq *proto.CommentRequest) (*proto.CommentListResponse, error) {
	list, err := g.VideoModel.GetComments(commentListReq.VideoID, commentListReq.CurrUserID)

	return &proto.CommentListResponse{
		Comments: list,
	}, err
}

func (g GRPCServer) GetVideoRecommendations(ctx context.Context, req *proto.RecReq) (*proto.RecResp, error) {
	resp, err := g.VideoModel.GetVideoRecommendations(req.UserId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
