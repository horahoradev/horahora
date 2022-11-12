package protocol

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/horahoradev/horahora/partyservice/internal/model"
	proto "github.com/horahoradev/horahora/partyservice/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type partyServer struct {
	proto.UnimplementedPartyserviceServer
	*model.PartyRepo
	mut      sync.Mutex
	timerMap map[int]*time.Timer
}

func New() (*partyServer, error) {

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
}

func (p *partyServer) GetPartyState(ctx context.Context, req *proto.PartyRequest) (*proto.PartyState, error) {
	return p.PartyRepo.GetPartyState(int(req.PartyID))
}
func (p *partyServer) NextVideo(ctx context.Context, req *proto.PartyRequest) (*proto.Empty, error) {
	// TODO: check: is leader??

	return &proto.Empty{}, p.PartyRepo.NextVideo(int(req.PartyID))
}
