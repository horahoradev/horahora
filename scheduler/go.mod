module github.com/horahoradev/horahora/scheduler

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-redsync/redsync v1.4.2
	github.com/go-stomp/stomp/v3 v3.0.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/horahoradev/horahora/video_service v0.0.0-20210922030328-dee03a489f47
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.4.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.14.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.28.1
)
