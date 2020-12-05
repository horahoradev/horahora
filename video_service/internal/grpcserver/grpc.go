package grpcserver

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/horahoradev/horahora/video_service/internal/dashutils"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/horahoradev/horahora/video_service/internal/models"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	proto "github.com/horahoradev/horahora/video_service/protocol"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	_ "github.com/aws/aws-sdk-go-v2/aws/external"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/google/uuid"
)

const uploadDir = "/videoservice/test_files/"

var _ proto.VideoServiceServer = (*GRPCServer)(nil)

type GRPCServer struct {
	VideoModel *models.VideoModel
	S3Client   s3.Client
	BucketName string
	Local      bool
	OriginFQDN string
}

func NewGRPCServer(bucketName string, db *sqlx.DB, port int, originFQDN string, local bool,
	redisClient *redis.Client, client userproto.UserServiceClient, tracer opentracing.Tracer) error {
	g, err := initGRPCServer(bucketName, db, client, local, redisClient, originFQDN)
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
	redisClient *redis.Client, originFQDN string) (*GRPCServer, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(cfg)

	g := &GRPCServer{
		S3Client:   *s3Client,
		BucketName: bucketName,
		Local:      local,
		OriginFQDN: originFQDN,
	}

	g.VideoModel, err = models.NewVideoModel(db, client, redisClient)
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

	tmpFile, err := ioutil.TempFile(uploadDir, id.String())
	if err != nil {
		err = fmt.Errorf("could not create tmp file. Err: %s", err)
		log.Error(err)
		return err
	}

	// This is awkward but it could be worse
	metaTmp, err := ioutil.TempFile(uploadDir, id.String()+".json")
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

	err = ioutil.WriteFile(video.FileData.Name()+".jpg", video.Meta.Meta.Thumbnail, 0644)
	if err != nil {
		return LogAndRetErr("could not write thumbnail. Err: %s", err)
	}

	log.Infof("Finished receiving file data for %s", video.Meta.Meta.Title)

	// If not local, upload the thumbnail and original video before returning
	if !g.Local {
		// FIXME did it again...
		err = g.SendToOriginServer(video.FileData.Name()+".jpg", filepath.Base(video.FileData.Name()+".jpg"))
		if err != nil {
			return err
		}

		// Upload the raw metadata
		err = g.SendToOriginServer(video.MetaFileData.Name(), filepath.Base(video.MetaFileData.Name()))
		if err != nil {
			return err
		}
		// Upload the original video
		err = g.SendToOriginServer(video.FileData.Name(), filepath.Base(video.FileData.Name()))
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
		video.Meta.Meta.AuthorUsername, video.Meta.Meta.AuthorUID, userproto.Site(video.Meta.Meta.OriginalSite),
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

				vid, err := g.fetchFromS3(video.GetMPDUUID())

				if err != nil {
					log.Errorf("could not fetch unencoded video id %d from s3. Err: %s", video.ID, err)
					return
				} else {
					defer func() {
						vid.Close()
						os.Remove(vid.Name())
					}()
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
	err := g.SendToOriginServer(d.ManifestPath, filepath.Base(d.ManifestPath))
	if err != nil {
		return err
	}

	// Send all of the chunked files
	for _, path := range d.QualityMap {
		err = g.SendToOriginServer(path, filepath.Base(path))
		if err != nil {
			return err
		}
	}

	return nil
}

func (g GRPCServer) fetchFromS3(id string) (*os.File, error) {
	getReq := &s3.GetObjectInput{
		Bucket: &g.BucketName,
		Key:    &id,
	}

	getObjReq := g.S3Client.GetObjectRequest(getReq)
	res, err := getObjReq.Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := os.OpenFile(uploadDir+id, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

loop:
	for true {
		buf := make([]byte, 1024*1024*1)
		n, err := res.Body.Read(buf)

		switch {
		case n == 0 && err == io.EOF:
			break loop
		case err != io.EOF && err != nil:
			err = fmt.Errorf("could not read from S3 for transcoding. Err: %s", err)
			log.Error(err)
			f.Close()
			return nil, err
		}

		// Truncate
		buf = buf[:n]

		_, err = f.Write(buf)
		if err != nil {
			f.Close()
			return nil, err
		}
	}

	return f, nil
}

func (g GRPCServer) SendToOriginServer(path, desiredFilename string) error {
	data, err := os.Open(path)
	if err != nil {
		return err
	}

	putObjInp := s3.PutObjectInput{
		ACL:                       s3.ObjectCannedACLPublicRead,
		Body:                      data,
		Bucket:                    &g.BucketName,
		CacheControl:              nil,
		ContentDisposition:        nil,
		ContentEncoding:           nil,
		ContentLanguage:           nil,
		ContentLength:             nil,
		ContentMD5:                nil,
		ContentType:               nil, // TODO
		Expires:                   nil,
		GrantFullControl:          nil,
		GrantRead:                 nil,
		GrantReadACP:              nil,
		GrantWriteACP:             nil,
		Key:                       &desiredFilename,
		Metadata:                  nil,
		ObjectLockLegalHoldStatus: "",
		ObjectLockMode:            "",
		ObjectLockRetainUntilDate: nil,
		RequestPayer:              "",
		SSECustomerAlgorithm:      nil,
		SSECustomerKey:            nil,
		SSECustomerKeyMD5:         nil,
		SSEKMSEncryptionContext:   nil,
		SSEKMSKeyId:               nil,
		ServerSideEncryption:      "",
		StorageClass:              "",
		Tagging:                   nil,
		WebsiteRedirectLocation:   nil,
	}

	putObjReq := g.S3Client.PutObjectRequest(&putObjInp)
	_, err = putObjReq.Send(context.TODO())
	return err // TODO
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
			queryConfig.FromUserID, queryConfig.ContainsTag, queryConfig.ShowUnapproved, queryConfig.OrderBy)
		if err != nil {
			log.Errorf("Could not get video list. Err: %s", err)
			return nil, err
		}

		numberOfVideos, err := g.VideoModel.GetNumberOfSearchResultsForQuery(queryConfig.FromUserID, queryConfig.ContainsTag)
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
