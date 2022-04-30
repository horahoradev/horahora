package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redsync/redsync"
	"github.com/horahoradev/horahora/scheduler/internal/models"

	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type schedulerServer struct {
	proto.UnimplementedSchedulerServer
	M *models.ArchiveRequestRepo
}

func NewGRPCServer(ctx context.Context, conn *sqlx.DB, rs *redsync.Redsync, port int) error {
	schedulerServer := initializeSchedulerServer(conn, rs)

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

func initializeSchedulerServer(conn *sqlx.DB, rs *redsync.Redsync) schedulerServer {
	return schedulerServer{
		M: models.NewArchiveRequest(conn, rs),
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

func (s schedulerServer) ListArchivalEntries(ctx context.Context, req *proto.ListArchivalEntriesRequest) (*proto.ListArchivalEntriesResponse, error) {
	archives, events, err := s.M.GetContentArchivalRequests(req.UserID)
	if err != nil {
		return nil, err
	}

	var entries []*proto.ContentArchivalEntry

	for _, archive := range archives {
		entry := proto.ContentArchivalEntry{
			UserID:             0, // In the future, will be expanded to allow queries for different users archival requests
			Url:                archive.Url,
			ArchivedVideos:     archive.Numerator,
			CurrentTotalVideos: archive.Denominator,
			BackoffFactor:      archive.BackoffFactor,
			LastSynced:         archive.LastSynced,
			DownloadID:         archive.DownloadID,
		}

		entries = append(entries, &entry)
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

	resp := proto.ListArchivalEntriesResponse{
		Entries: entries,
		Events:  protoEvents,
	}

	return &resp, nil
}
