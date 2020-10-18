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
	Meta     *proto.InputVideoChunk_Meta
	FileData *os.File
}

func (g GRPCServer) UploadVideo(inpStream proto.VideoService_UploadVideoServer) error {
	log.Info("Handling video upload")

	var video VideoUpload

	uploadDir := "/videoservice/test_files/"

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

	video.FileData = tmpFile
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

	transcodeResults, err := dashutils.TranscodeAndGenerateManifest(video.FileData.Name(), g.Local)
	if err != nil {
		err := fmt.Errorf("failed to transcode and chunk. Err: %s", err)
		log.Error(err)
		return err
	}

	// S3? Nginx? Who knows...
	if !g.Local {
		err = g.UploadMPDSet(transcodeResults)
		if err != nil {
			err := fmt.Errorf("failed to upload to origin. Err: %s", err)
			log.Error(err)
			return err
		}
	}

	// This is MESSY
	// TODO: switch to struct for args
	manifestLoc := fmt.Sprintf("%s/%s", g.OriginFQDN, filepath.Base(transcodeResults.ManifestPath))

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

	// Upload thumbnail
	// FIXME did it again...
	err = g.SendToOriginServer(d.ThumbnailPath, filepath.Base(d.ThumbnailPath))
	if err != nil {
		return err
	}

	// Upload the original video
	err = g.SendToOriginServer(d.OriginalFilePath, filepath.Base(d.OriginalFilePath))
	if err != nil {
		return err
	}

	return nil
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
	list, err := g.VideoModel.GetComments(commentListReq.VideoID)

	return &proto.CommentListResponse{
		Comments: list,
	}, err
}
