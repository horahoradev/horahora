module github.com/horahoradev/horahora/frontend

go 1.13

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/horahoradev/horahora/scheduler v0.0.0-20200823010907-db67efcf7d8c
	github.com/horahoradev/horahora/user_service v0.0.0-20200929200329-a2cc6bce4184
	github.com/horahoradev/horahora/video_service v0.0.0-20201021061128-99df23c0d6e6
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/moxiaomomo/grpc-jaeger v0.0.0-20180617090213-05b879580c4a // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.5.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	google.golang.org/grpc v1.33.0
)
