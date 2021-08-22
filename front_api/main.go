package main

import (
	"github.com/horahoradev/horahora/front_api/config"
	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	routes "github.com/horahoradev/horahora/front_api/routes"

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

	grpcAuth := custommiddleware.NewGRPCAuth(cfg)
	e.Use(grpcAuth.GRPCAuth)

	routes.SetupRoutes(e, cfg)

	e.Logger.Fatal(e.Start(":8083"))
}
