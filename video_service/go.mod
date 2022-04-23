module github.com/horahoradev/horahora/video_service

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.16.2
	github.com/aws/aws-sdk-go-v2/config v1.15.3
	github.com/aws/aws-sdk-go-v2/credentials v1.11.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.5
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/doug-martin/goqu/v9 v9.9.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/mock v1.4.3
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/horahoradev/horahora/user_service v0.0.0-20210922030328-dee03a489f47
	github.com/jmoiron/sqlx v1.2.0
	github.com/kurin/blazer v0.5.3
	github.com/lib/pq v1.4.0
	github.com/minio/minio-go/v7 v7.0.10
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
