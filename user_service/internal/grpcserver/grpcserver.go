package grpcserver

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"

	"github.com/horahoradev/horahora/user_service/internal/auth"
	"github.com/horahoradev/horahora/user_service/internal/model"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"

	proto "github.com/horahoradev/horahora/user_service/protocol"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	proto.UnsafeUserServiceServer
	db         *sqlx.DB
	privateKey *rsa.PrivateKey
	um         *model.UserModel
}

// Compile-time implementation check
var _ proto.UserServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(db *sqlx.DB, privateKey *rsa.PrivateKey, port int64) error {
	um, err := model.NewUserModel(db)
	if err != nil {
		return err
	}

	g := GRPCServer{
		db:         db,
		privateKey: privateKey,
		um:         um,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Infof("Listening on port %d", port)
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, g)
	return grpcServer.Serve(lis)
}

func (g GRPCServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if !req.ForeignUser {
		if err := req.Validate(); err != nil {
			return nil, err
		}
	}

	log.Infof("Handling registration for user %s", req.Username)
	jwt, err := auth.Register(req.Username, req.Email, req.Password, g.um, g.privateKey, req.ForeignUser, req.ForeignUserID, req.ForeignWebsite)
	if err != nil {
		log.Errorf("auth: failed to register user %s, failed with err %s", req.Username, err)
		return nil, err
	}

	p := proto.RegisterResponse{
		Jwt: jwt,
	}

	return &p, nil
}

func (g GRPCServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	log.Info("Handling login for user %s", req.Username)
	jwt, err := auth.Login(req.Username, req.Password, g.privateKey, g.um)
	if err != nil {
		log.Errorf("auth login failed with err: %s", err)
		return nil, err
	}

	p := proto.LoginResponse{
		Jwt: jwt,
	}

	return &p, nil
}

func (g GRPCServer) BanUser(ctx context.Context, req *proto.BanUserRequest) (*proto.BanUserResponse, error) {
	idToBan := req.UserID
	err := g.um.BanUser(idToBan)
	return &proto.BanUserResponse{}, err
}

func (g GRPCServer) GetUserFromID(ctx context.Context, req *proto.GetUserFromIDRequest) (*proto.UserResponse, error) {
	id := req.UserID

	user, err := g.um.GetUserWithID(id)
	if err != nil {
		log.Errorf("failed to fetch user with id %s, failed with err %s", id, err)
		return nil, err
	}

	return &proto.UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Rank:     proto.UserRank(user.Rank),
		Banned:   user.Banned,
	}, nil
}

func (g GRPCServer) ValidateJWT(ctx context.Context, req *proto.ValidateJWTRequest) (*proto.ValidateJWTResponse, error) {
	uid, err := auth.ValidateJWT(req.Jwt, *g.privateKey)
	if err != nil {
		log.Errorf("failed to validate JWT, err: %s", err)
		return nil, err
	}

	banned, err := g.um.IsBanned(uid)
	if err != nil {
		return nil, err
	}

	if banned {
		return nil, errors.New("User is banned")
	}

	return &proto.ValidateJWTResponse{
		IsValid: true,
		Uid:     uid,
	}, nil
}

func (g GRPCServer) GetUserIDsForUsername(ctx context.Context, req *proto.GetUserIDsForUsernameRequest) (*proto.GetUserIDsForUsernameResponse, error) {
	ids, err := g.um.GetUserIDsForUsername(req.Username)
	if err != nil {
		log.Errorf("Failed to retrieve usernames for id. Err: %s", err)
		return nil, err
	}

	return &proto.GetUserIDsForUsernameResponse{
		UserIDs: ids,
	}, nil
}

func (g GRPCServer) GetUserForForeignUID(ctx context.Context, req *proto.GetForeignUserRequest) (*proto.GetForeignUserResponse, error) {
	uid, err := g.um.GetForeignUser(req.ForeignUserID, req.OriginalWebsite)
	if err != nil {
		log.Errorf("failed to get foreign UID, err: %s", err)
		return nil, err
	}

	return &proto.GetForeignUserResponse{
		NewUID: uid,
	}, nil
}

func (g GRPCServer) SetUserRank(ctx context.Context, req *proto.SetRankRequest) (*proto.Nothing, error) {
	return &proto.Nothing{}, g.um.SetUserRank(req.UserID, int64(req.Rank.Number()))
}

func (g GRPCServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.Nothing, error) {
	// This is a little inefficient but whatever
	user, err := g.um.GetUserWithID(req.UserID)
	if err != nil {
		log.Errorf("failed to fetch user with id %s, failed with err %s", req.UserID, err)
		return nil, err
	}

	_, err = auth.Login(user.Username, req.OldPassword, g.privateKey, g.um)
	if err != nil {
		log.Errorf("Password reset auth failed with err: %s", err)
		return nil, err
	}

	// old password is valid, so we can proceed to creating a new hash
	passHash, err := auth.GenerateHash([]byte(req.NewPassword))
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, g.um.SetNewHash(req.UserID, passHash)
}

func (g GRPCServer) AddAuditEvent(ctx context.Context, req *proto.NewAuditEventRequest) (*proto.Nothing, error) {
	return &proto.Nothing{}, g.um.AddNewAuditEvent(req.User_ID, req.Message)
}

func (g GRPCServer) GetAuditEvents(ctx context.Context, req *proto.AuditEventsListRequest) (*proto.AuditListResponse, error) {
	events, count, err := g.um.GetAuditEvents(req.UserId, req.Page)
	if err != nil {
		return nil, err
	}
	auditEvents := make([]*proto.AuditEvent, 0, len(events))

	for _, event := range events {
		event := proto.AuditEvent{
			Id:        event.ID,
			Message:   event.Message,
			Timestamp: event.Timestamp,
			User_ID:   event.UserID,
		}

		auditEvents = append(auditEvents, &event)
	}

	return &proto.AuditListResponse{
		Events:    auditEvents,
		NumEvents: count,
	}, nil
}
