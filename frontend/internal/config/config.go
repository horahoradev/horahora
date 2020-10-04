package config

import (
	"github.com/caarlos0/env"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"google.golang.org/grpc"

	"time"
)

type Config struct {
	UserServiceGRPCAddress      string `env:"UserServiceGRPCAddress",required"`
	VideoServiceGRPCAddress     string `env:"VideoServiceGRPCAddress",required"`
	SchedulerServiceGRPCAddress string `env:"SchedulerServiceGRPCAddress",required"`

	VideoClient     videoproto.VideoServiceClient
	UserClient      userproto.UserServiceClient
	SchedulerClient schedulerproto.SchedulerClient
}

func New() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	videoGRPCConn, err := grpc.Dial(config.VideoServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		return nil, err
	}

	userGRPCConn, err := grpc.Dial(config.UserServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		return nil, err
	}

	schedulerGRPCConn, err := grpc.Dial(config.SchedulerServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		return nil, err
	}

	config.SchedulerClient = schedulerproto.NewSchedulerClient(schedulerGRPCConn)
	config.UserClient = userproto.NewUserServiceClient(userGRPCConn)
	config.VideoClient = videoproto.NewVideoServiceClient(videoGRPCConn)

	return &config, nil
}
