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

	err := s.M.New(models.Channel, string(req.ChannelID), req.Website)

	return ret, err
}

func (s schedulerServer) DlPlaylist(ctx context.Context, req *proto.PlaylistRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(models.Playlist, req.PlaylistID, req.Website)

	return ret, err
}

func (s schedulerServer) DlTag(ctx context.Context, req *proto.TagRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	err := s.M.New(models.Tag, req.TagValue, req.Website)

	return ret, err
}
