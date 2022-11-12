package protocol

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/horahoradev/horahora/partyservice/internal/model"
	proto "github.com/horahoradev/horahora/partyservice/protocol"
	videoservice "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type partyServer struct {
	proto.UnimplementedPartyserviceServer
	PartyRepo *model.PartyRepo
	mut       sync.Mutex
	timerMap  map[int]*time.Timer
	v         videoservice.VideoServiceClient
}

func New(db *sqlx.DB, v videoservice.VideoServiceClient) (*partyServer, error) {
	return &partyServer{
		PartyRepo: model.New(db),
		mut:       sync.Mutex{},
		timerMap:  make(map[int]*time.Timer),
		v:         v,
	}, nil
}

func (p *partyServer) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPartyserviceServer(grpcServer, p)
	return grpcServer.Serve(lis)

}

func (p *partyServer) BecomeLeader(context.Context, *proto.PartyRequest) (*proto.LeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BecomeLeader not implemented")
}

func (p *partyServer) JoinParty(ctx context.Context, req *proto.PartyRequest) (*proto.Empty, error) {
	err := p.PartyRepo.JoinWatchParty(int(req.UserID), int(req.PartyID))
	if err != nil {
		return nil, err
	}

	// lmfao
	timer := time.NewTimer(time.Second * 30)

	p.mut.Lock()
	p.timerMap[int(req.UserID)] = timer
	p.mut.Unlock()

	go func() {
		<-timer.C
		p.mut.Lock()
		defer p.mut.Unlock()

		p.PartyRepo.DeleteFromWatchParty(int(req.UserID), int(req.PartyID))
		// lol idk what to do with this error
		// TODO log me
		delete(p.timerMap, int(req.UserID))
	}()
	return &proto.Empty{}, nil
}

func (p *partyServer) HeartBeat(ctx context.Context, req *proto.PartyRequest) (*proto.Empty, error) {
	p.mut.Lock()
	timer, ok := p.timerMap[int(req.UserID)]
	defer p.mut.Unlock()
	if !ok {
		return nil, errors.New("timer not found in map")
	}

	// Back to 0
	timer.Reset(time.Second * 30)
	return &proto.Empty{}, nil
}

func (p *partyServer) GetPartyState(ctx context.Context, req *proto.PartyRequest) (*proto.PartyState, error) {
	return p.PartyRepo.GetPartyState(int(req.PartyID))
}

func (p *partyServer) AddVideo(ctx context.Context, req *proto.VideoRequest) (*proto.Empty, error) {
	url, err := url.Parse(req.VideoURL)
	if err != nil {
		return nil, err
	}

	// happy path
	// FIXME
	videoID := url.Path[len(url.Path)-1]

	resp, err := p.v.GetVideo(context.Background(), &videoservice.VideoRequest{
		VideoID: fmt.Sprintf("%v", videoID),
	})
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, p.PartyRepo.NewVideo(resp.VideoLoc, resp.VideoTitle, int(resp.VideoID))
}

func (p *partyServer) NextVideo(ctx context.Context, req *proto.PartyRequest) (*proto.Empty, error) {
	// TODO: check: is leader??

	return &proto.Empty{}, p.PartyRepo.NextVideo(int(req.PartyID))
}
