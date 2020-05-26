package grpcserver

import (
	"context"
	"fmt"
	"net"

	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type schedulerServer struct {
	proto.UnimplementedSchedulerServer
	Db *sqlx.DB
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
		Db: conn,
	}
}

func (s schedulerServer) DlChannel(ctx context.Context, req *proto.ChannelRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	website, err := getWebsiteStringFromEnum(req.Website)
	if err != nil {
		return ret, err
	}

	_, err = s.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value) "+
		"VALUES (Now(), $1, $2, $3)", website, "channel", req.ChannelID)

	return ret, err
}

func (s schedulerServer) DlPlaylist(ctx context.Context, req *proto.PlaylistRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	website, err := getWebsiteStringFromEnum(req.Website)
	if err != nil {
		return ret, err
	}

	// TODO: implement num to download
	_, err = s.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value) "+
		"VALUES (Now(), $1, $2, $3)", website, "playlist", req.PlaylistID)

	return ret, err
}
func (s schedulerServer) DlTag(ctx context.Context, req *proto.TagRequest) (*proto.Empty, error) {
	ret := &proto.Empty{}

	website, err := getWebsiteStringFromEnum(req.Website)
	if err != nil {
		return ret, err
	}

	_, err = s.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value) "+
		"VALUES (Now(), $1, $2, $3)", website, "tag", req.TagValue)

	return ret, err
}

// FIXME: this is dumb
func getWebsiteStringFromEnum(enumVal proto.Site) (string, error) {
	switch enumVal {
	case proto.Site_niconico:
		return "niconico", nil
	case proto.Site_bilibili:
		return "bilibili", nil
	case proto.Site_youtube:
		return "youtube", nil
	default:
		return "", fmt.Errorf("could not find specified website for enum %d", enumVal)
	}
}
