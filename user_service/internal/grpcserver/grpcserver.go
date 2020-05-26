package grpcserver

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"

	"git.horahora.org/otoman/user-service.git/internal/auth"
	"git.horahora.org/otoman/user-service.git/internal/model"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"

	proto "git.horahora.org/otoman/user-service.git/protocol"
	"google.golang.org/grpc"
)

type GRPCServer struct {
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
	}, nil
}

func (g GRPCServer) ValidateJWT(ctx context.Context, req *proto.ValidateJWTRequest) (*proto.ValidateJWTResponse, error) {
	uid, err := auth.ValidateJWT(req.Jwt, *g.privateKey)
	if err != nil {
		log.Errorf("failed to validate JWT, err: %s", err)
		return nil, err
	}

	return &proto.ValidateJWTResponse{
		IsValid: true,
		Uid:     uid,
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
