package main

import (
	"github.com/horahoradev/horahora/frontend/internal/config"
	"net/http"

	custommiddleware "github.com/horahoradev/horahora/frontend/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config. Err: %s", err)
	}

	e := echo.New()

	grpcAuth := custommiddleware.NewGRPCAuth(cfg)

	e.Use(middleware.Logger())
	e.Use(grpcAuth.GRPCAuth)

	setupRoutes(e)

	e.Logger.Fatal(e.Start(":80"))
}

func setupRoutes(e *echo.Echo) {
	e.GET("/users/:id", getUser)
	e.GET("/videos/:id", getVideo)
	e.GET("/", getHome)
}

// This section is just a placeholder

func getUser(c echo.Context) error {
	// TODO
	// User ID from path `users/:id`
	id := c.Param("id")

	return c.String(http.StatusOK, id)
}

func getVideo(c echo.Context) error {
	// TODO
	id := c.Param("id")

	return c.String(http.StatusOK, id)
}

func getHome(c echo.Context) error {
	// TODO
	return c.String(http.StatusOK, "home")
}
