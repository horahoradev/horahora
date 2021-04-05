package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"net"
	"time"
)

type PostgresInfo struct {
	Hostname string `env:"pgs_host,required"`
	Port     int    `env:"pgs_port,required"`
	Username string `env:"pgs_user"`
	Password string `env:"pgs_pass"`
	Db       string `env:"pgs_db,required"`
}

type RedisInfo struct {
	Hostname string `env:"redis_host,required"`
	Port     int    `env:"redis_port,required"`
	Password string `env:"redis_pass,required"`
}

type config struct {
	PostgresInfo
	RedisInfo
	RedisConn              *redis.Client
	GRPCPort               int    `env:"GRPCPort,required"`
	UserServiceGRPCAddress string `env:"UserServiceGRPCAddress,required"`
	BucketName             string `env:"BucketName,required"`
	Local                  bool   `env:"Local,required"` // If running locally, no s3 uploads
	// (this is a workaround for getting IAM permissions into pods running on minikube)
	OriginFQDN        string `env:"OriginFQDN,required"`
	JaegerAddress     string `env:"JaegerAddress"`
	UserClient        userproto.UserServiceClient
	SqlClient         *sqlx.DB
	Tracer            opentracing.Tracer
	StorageBackend    string `env:"StorageBackend,required"`
	StorageAPIID      string `env:"StorageAPIID"`
	StorageAPIKey     string `env:"StorageAPIKey"`
	ApprovalThreshold int    `env:"ApprovalThreshold,required"`
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.RedisInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}

	config.RedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisInfo.Hostname, config.RedisInfo.Port),
		Password: config.RedisInfo.Password, // no password set
		DB:       0,                         // use default DB
	})

	config.SqlClient, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.PostgresInfo.Hostname, config.PostgresInfo.Username, config.PostgresInfo.Password, config.PostgresInfo.Db))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to postgres. Err: %s", err)
	}

	transport, err := jaeger.NewUDPTransport(net.JoinHostPort(config.JaegerAddress, "6831"),
		4096)
	if err != nil {
		log.Errorf("Could not initialize jaeger tracer. Err: %s", err)
	}

	tracer, _ := jaeger.NewTracer("videoservice",
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
	)

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	conn, err := grpc.Dial(config.UserServiceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			otgrpc.OpenTracingClientInterceptor(tracer),
			grpc_retry.UnaryClientInterceptor(opts...))))
	if err != nil {
		return nil, err
	}

	config.UserClient = userproto.NewUserServiceClient(conn)
	config.Tracer = tracer

	return &config, err
}
