package grpcserver

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/horahoradev/horahora/scheduler/internal/models"
	"net"

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

func (s schedulerServer) ListArchivalEntries(ctx context.Context, req *proto.ListArchivalEntriesRequest) (*proto.ListArchivalEntriesResponse, error) {
	urls, err := s.M.GetContentArchivalRequests(req.UserID)
	if err != nil {
		return nil, err
	}

	var entries []*proto.ContentArchivalEntry

	for _, url := range urls {
		entry := proto.ContentArchivalEntry{
			UserID:       0, // In the future, will be expanded to allow queries for different users archival requests
			Url: url,
		}

		entries = append(entries, &entry)
	}

	resp := proto.ListArchivalEntriesResponse{
		Entries: entries,
	}
	return &resp, nil
}
