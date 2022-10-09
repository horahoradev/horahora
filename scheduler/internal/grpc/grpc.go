package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/horahoradev/horahora/scheduler/internal/models"

	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type schedulerServer struct {
	proto.UnimplementedSchedulerServer
	M *models.ArchiveRequestRepo
}

func NewGRPCServer(ctx context.Context, conn *sqlx.DB, port int) error {
	schedulerServer := initializeSchedulerServer(conn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	proto.RegisterSchedulerServer(serv, schedulerServer)

	go func() {
		<-ctx.Done()
		serv.GracefulStop()
	}()

	return serv.Serve(lis)
}

func initializeSchedulerServer(conn *sqlx.DB) schedulerServer {
	return schedulerServer{
		M: models.NewArchiveRequest(conn),
	}
}

func (s schedulerServer) DlURL(ctx context.Context, req *proto.URLRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(req.Url, req.UserID)

	return ret, err
}

func (s schedulerServer) DeleteArchivalRequest(ctx context.Context, req *proto.DeletionRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.DeleteArchivalRequest(req.UserID, req.DownloadID)

	return ret, err
}

func (s schedulerServer) ListArchivalEvents(ctx context.Context, req *proto.ListArchivalEventsRequest) (*proto.ListArchivalEventsResponse, error) {
	events, err := s.M.GetArchivalEvents(req.DownloadID, req.ShowAll)
	if err != nil {
		return nil, err
	}

	var protoEvents []*proto.ArchivalEvent
	for _, event := range events {
		eventObj := proto.ArchivalEvent{
			VideoUrl:  event.VideoURL,
			ParentUrl: event.ParentURL,
			Message:   event.Message,
			Timestamp: event.EventTimestamp,
		}
		protoEvents = append(protoEvents, &eventObj)
	}

	return &proto.ListArchivalEventsResponse{
		Events: protoEvents,
	}, nil
}

func (s schedulerServer) ListArchivalEntries(ctx context.Context, req *proto.ListArchivalEntriesRequest) (*proto.ListArchivalEntriesResponse, error) {
	archives, err := s.M.GetContentArchivalRequests(req.UserID)
	if err != nil {
		return nil, err
	}

	var entries []*proto.ContentArchivalEntry

	for _, archive := range archives {
		entry := proto.ContentArchivalEntry{
			UserID:               0, // In the future, will be expanded to allow queries for different users archival requests
			Url:                  archive.Url,
			ArchivedVideos:       archive.Numerator,
			CurrentTotalVideos:   archive.Denominator,
			BackoffFactor:        archive.BackoffFactor,
			LastSynced:           archive.LastSynced,
			DownloadID:           archive.DownloadID,
			UndownloadableVideos: archive.Undownloadable,
		}

		entries = append(entries, &entry)
	}

	resp := proto.ListArchivalEntriesResponse{
		Entries: entries,
	}

	return &resp, nil
}

func (s schedulerServer) RetryArchivalRequestDownloadss(ctx context.Context, req *proto.RetryRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.RetryArchivalRequest(req.UserID, req.DownloadID)

	return ret, err
}

func (s schedulerServer) GetDownloadsInProgress(ctx context.Context, req *proto.DownloadsInProgressRequest) (*proto.DownloadsInProgressResponse, error) {
	videos, err := s.M.GetDownloadsInProgress()
	if err != nil {
		return nil, err
	}

	var ret []*proto.Video
	for _, video := range videos {
		vid := proto.Video{
			VideoID: video.ID,
			Website: video.Website,
		}
		if video.DlStatus == 3 {
			// downloading
			vid.DlStatus = proto.Video_Downloading
		} else if video.DlStatus == 4 {
			// queued
			vid.DlStatus = proto.Video_Queued
		}
		ret = append(ret, &vid)
	}

	return &proto.DownloadsInProgressResponse{Videos: ret}, nil
}
