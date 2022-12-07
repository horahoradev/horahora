package main

import (
	"github.com/horahoradev/horahora/front_api/config"
	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	routes "github.com/horahoradev/horahora/front_api/routes"
	"github.com/horahoradev/horahora/front_api/sockets"
	"github.com/labstack/echo-contrib/prometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config. Err: %s", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	grpcAuth := custommiddleware.NewGRPCAuth(cfg)
	e.Use(grpcAuth.GRPCAuth)

	srv := sockets.New()

	routes.SetupRoutes(e, cfg, srv)

	go sockets.Run(srv)

	e.Logger.Fatal(e.Start(":8083"))
}
