package grpcserver

import (
	"context"
	"fmt"
	"github.com/horahoradev/horahora/scheduler/internal/models"
	"net"

	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type schedulerServer struct {
	proto.UnimplementedSchedulerServer
	M *models.ArchiveRequestModel
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

func (s schedulerServer) DlChannel(ctx context.Context, req *proto.ChannelRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(models.Channel, string(req.ChannelID), req.Website, req.UserID)

	return ret, err
}

func (s schedulerServer) DlPlaylist(ctx context.Context, req *proto.PlaylistRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(models.Playlist, req.PlaylistID, req.Website, req.UserID)

	return ret, err
}

func (s schedulerServer) DlTag(ctx context.Context, req *proto.TagRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(models.Tag, req.TagValue, req.Website, req.UserID)

	return ret, err
}

func (s schedulerServer) ListArchivalEntries(ctx context.Context, req *proto.ListArchivalEntriesRequest) (*proto.ListArchivalEntriesResponse, error) {
	requests, err := s.M.GetContentArchivalRequests(req.UserID)
	if err != nil {
		return nil, err
	}

	var entries []*proto.ContentArchivalEntry

	for _, request := range requests {
		entry := proto.ContentArchivalEntry{
			UserID:       0, // In the future, will be expanded to allow queries for different users archival requests
			Website:      request.Website,
			ContentType:  string(request.ContentType),
			ContentValue: request.ContentValue,
		}

		entries = append(entries, &entry)
	}

	resp := proto.ListArchivalEntriesResponse{
		Entries: entries,
	}
	return &resp, nil
}
