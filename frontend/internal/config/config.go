package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

type Config struct {
	UserServiceGRPCAddress      string `env:"UserServiceGRPCAddress,required"`
	VideoServiceGRPCAddress     string `env:"VideoServiceGRPCAddress,required"`
	SchedulerServiceGRPCAddress string `env:"SchedulerServiceGRPCAddress,required"`

	JaegerAddress string `env:"JaegerAddress,required"`

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

	// 4096 is an okay max packet size I guess?
	transport, err := jaeger.NewUDPTransport(fmt.Sprintf("%s:6832", config.JaegerAddress), 4096)
	if err != nil {
		return nil, err
	}

	// TODO: close when main exits to flush traces
	tracer, _ := jaeger.NewTracer("frontend",
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
		jaeger.TracerOptions.Logger(jaeger.StdLogger),
	)

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	// TODO: am I doing this right?
	videoGRPCConn, err := grpc.Dial(config.VideoServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			otgrpc.OpenTracingClientInterceptor(tracer),
			grpc_retry.UnaryClientInterceptor(opts...)),
		))
	if err != nil {
		return nil, err
	}

	userGRPCConn, err := grpc.Dial(config.UserServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			otgrpc.OpenTracingClientInterceptor(tracer),
			grpc_retry.UnaryClientInterceptor(opts...)),
		))
	if err != nil {
		return nil, err
	}

	schedulerGRPCConn, err := grpc.Dial(config.SchedulerServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			otgrpc.OpenTracingClientInterceptor(tracer),
			grpc_retry.UnaryClientInterceptor(opts...)),
		))
	if err != nil {
		return nil, err
	}

	config.SchedulerClient = schedulerproto.NewSchedulerClient(schedulerGRPCConn)
	config.UserClient = userproto.NewUserServiceClient(userGRPCConn)
	config.VideoClient = videoproto.NewVideoServiceClient(videoGRPCConn)

	return &config, nil
}
