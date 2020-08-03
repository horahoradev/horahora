package main

import (
	"github.com/horahoradev/horahora/frontend/internal/config"
	"github.com/horahoradev/horahora/frontend/internal/templates"
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

	t := templates.New()

	e.Renderer = t

	e.Logger.Fatal(e.Start(":80"))
}

func setupRoutes(e *echo.Echo) {
	e.GET("/users/:id", getUser)
	e.GET("/videos/:id", getVideo)
	e.GET("/", getHome)
}

type Video struct {
	Title    string
	Views    int64
	AuthorID int64
	Rating   float32
}

type ProfileData struct {
	ActiveUserID    int64
	ProfileUserID   int64
	ProfileUsername string
	ProfileVideos   []Video
}

func getUser(c echo.Context) error {
	// TODO
	// User ID from path `users/:id`
	//id := c.Param("id")

	data := ProfileData{}

	id, ok := c.Get(custommiddleware.UserIDKey).(*int64)
	if ok && id != nil {
		data.ActiveUserID = *id
	}
	// TODO: get the rest of data

	return c.Render(http.StatusOK, "profile", data)
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
