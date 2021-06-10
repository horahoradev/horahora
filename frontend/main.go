package main

import (
	"net/url"

	"github.com/horahoradev/horahora/frontend/internal/config"
	custommiddleware "github.com/horahoradev/horahora/frontend/internal/middleware"
	"github.com/horahoradev/horahora/frontend/internal/routes"
	"github.com/horahoradev/horahora/frontend/internal/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const testData = false

func main() {
	e := echo.New()
	e.Static("/static", "assets")
	e.Use(middleware.Logger())

	switch testData {
	case true:
		routes.SetupTestRoutes(e)
	case false:
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("Could not initialize config. Err: %s", err)
		}

		grpcAuth := custommiddleware.NewGRPCAuth(cfg)
		e.Use(grpcAuth.GRPCAuth)

		url1, err := url.Parse("http://nginx:86")
		if err != nil {
			e.Logger.Fatal(err)
		}

		targets := []*middleware.ProxyTarget{
			{
				URL: url1,
			},
		}
		g := e.Group("/staticfiles")
		g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
		routes.SetupRoutes(e, cfg)
	}

	t := templates.New()

	e.Renderer = t

	e.Logger.Fatal(e.Start(":8082"))
}
