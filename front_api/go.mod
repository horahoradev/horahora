module github.com/SEAPUNK/horahora/front_api

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/horahoradev/horahora/scheduler v0.0.0-00010101000000-000000000000
	github.com/horahoradev/horahora/user_service v0.0.0-20200526031340-64e1705d00d7
	github.com/horahoradev/horahora/video_service v0.0.0-20201205215129-690cef6cbea9
	github.com/labstack/echo/v4 v4.3.0
	github.com/labstack/gommon v0.3.0
	google.golang.org/grpc v1.38.0
)

replace github.com/horahoradev/horahora/scheduler => ../scheduler

replace github.com/horahoradev/horahora/user_service => ../user_service

replace github.com/horahoradev/horahora/video_service => ../video_service
