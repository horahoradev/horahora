package main

import (
	"context"
	"github.com/horahoradev/horahora/frontend/internal/config"
	custommiddleware "github.com/horahoradev/horahora/frontend/internal/middleware"
	"github.com/horahoradev/horahora/frontend/internal/templates"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config. Err: %s", err)
	}

	e := echo.New()

	grpcAuth := custommiddleware.NewGRPCAuth(cfg)
	e.Use(grpcAuth.GRPCAuth)
	e.Static("/static", "assets")
	e.Use(middleware.Logger())

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

	setupRoutes(e, cfg)

	t := templates.New()

	e.Renderer = t

	e.Logger.Fatal(e.Start(":8081"))
}

func setupRoutes(e *echo.Echo, cfg *config.Config) {

	h := NewHomeHandler(cfg.VideoClient)
	e.GET("/", h.getHome)
	e.GET("/users/:id", getUser)

	v := NewVideoHandler(cfg.VideoClient, cfg.UserClient)
	e.GET("/videos/:id", v.getVideo)
	e.GET("/login", getLogin)
	e.POST("/login", handleLogin)

	e.GET("/register", getRegister)
	e.POST("/register", handleRegistration)
}

type Video struct {
	Title        string
	VideoID      int64
	Views        uint64
	AuthorID     int64
	AuthorName   string
	ThumbnailLoc string
	Rating       float64
}

type Comment struct {
	ProfilePicture string
	Username       string
	Comment        string
}

type VideoDetail struct {
	L               LoggedInUserData
	Title           string
	MPDLoc          string
	Views           uint64
	Rating          float64
	AuthorID        int64
	Username        string
	UserDescription string
	UserSubscribers uint64
	ProfilePicture  string
	UploadDate      string // should be a datetime
	Comments        []Comment
}

type LoggedInUserData struct {
	UserID            int64
	Username          string
	ProfilePictureURL string
}

type ProfileData struct {
	L                 LoggedInUserData
	UserID            int64
	Username          string
	ProfilePictureURL string
	Videos            []Video
}
type HomePageData struct {
	L      LoggedInUserData
	Videos []Video
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func getRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func handleLogin(c echo.Context) error {
	//username := c.FormValue("username")
	//password := c.FormValue("password")

	// TODO: grpc auth goes here

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = "jwt" // JWT will go here after grpc logins are added
	//cookie.Expires = time.Now().Add(24 * time.Hour)

	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = true

	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Login successful.")

}

func handleRegistration(c echo.Context) error {
	//username := c.FormValue("username")
	//email := c.FormValue("email")
	//password := c.FormValue("password")

	return c.String(http.StatusOK, "Registration successful.")
}

func getUser(c echo.Context) error {
	// TODO
	// User ID from path `users/:id`
	//id := c.Param("id")

	// This is just sample data to make sure everything renders correctly
	data := ProfileData{
		Username:          "testuser",
		UserID:            1,
		ProfilePictureURL: "/static/images/placeholder1.jpg",
		Videos: []Video{
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
			{
				Title:        "testvideo",
				Views:        5,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
		},
	}

	id, ok := c.Get(custommiddleware.UserIDKey).(*int64)
	if ok && id != nil {
		data.L.UserID = *id
	}
	// TODO: get the rest of data

	return c.Render(http.StatusOK, "profile", data)
}

type VideoHandler struct {
	v videoproto.VideoServiceClient
	u userproto.UserServiceClient
}

func NewVideoHandler(v videoproto.VideoServiceClient, u userproto.UserServiceClient) *VideoHandler {
	return &VideoHandler{
		v: v,
		u: u,
	}
}

func (v *VideoHandler) getVideo(c echo.Context) error {
	id := c.Param("id")

	// Dumb
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// Increment views first
	viewReq := videoproto.VideoViewing{VideoID: idInt}
	_, err = v.v.ViewVideo(context.Background(), &viewReq)
	if err != nil {
		return err
	}

	videoReq := videoproto.VideoRequest{
		VideoID: id,
	}

	videoInfo, err := v.v.GetVideo(context.Background(), &videoReq)
	if err != nil {
		return err
	}

	spl := strings.Split(videoInfo.VideoLoc, "/")

	rating := videoInfo.Rating

	// lol
	if math.IsNaN(rating) {
		rating = 0.00
	}

	data := VideoDetail{
		L:               LoggedInUserData{},
		Title:           videoInfo.VideoTitle,
		MPDLoc:          spl[len(spl)-1], // FIXME: fix this in videoservice LOL this is embarrassing
		Views:           videoInfo.Views,
		Rating:          rating,
		AuthorID:        videoInfo.AuthorID, // TODO
		Username:        videoInfo.AuthorName,
		UserDescription: "", // TODO: not implemented yet
		UserSubscribers: 0,  // TODO: not implemented yet
		ProfilePicture:  "/static/images/placeholder1.jpg",
		UploadDate:      videoInfo.UploadDate,
		Comments:        nil,
	}

	//data := VideoDetail{
	//	Title:           "My cool video",
	//	MPDLoc:          "",
	//	Views:           100,
	//	Rating:          10.0,
	//	AuthorID:        4,
	//	Username:        "testuser",
	//	UserDescription: "we did it reddit",
	//	ProfilePicture:  "/static/images/placeholder1.jpg",
	//	UploadDate:      time.Now(),
	//	UserSubscribers: 100,
	//	Comments: []Comment{
	//		{
	//			ProfilePicture: "/static/images/placeholder1.jpg",
	//			Username:       "testuser2",
	//			Comment:        "WOW",
	//		},
	//	},
	//}

	return c.Render(http.StatusOK, "video", data)
}

type HomeHandler struct {
	videoClient videoproto.VideoServiceClient
}

func NewHomeHandler(v videoproto.VideoServiceClient) HomeHandler {
	return HomeHandler{videoClient: v}
}

func (h *HomeHandler) getHome(c echo.Context) error {
	// TODO: if request times out, maybe provide a default list of good videos
	req := videoproto.VideoQueryConfig{
		OrderBy:    videoproto.OrderCategory_upload_date,
		Direction:  videoproto.SortDirection_desc,
		PageNumber: 0,
	}

	videoList, err := h.videoClient.GetVideoList(context.TODO(), &req)
	if err != nil {
		log.Errorf("Could not retrieve video list. Err: %s", err)
		return c.String(http.StatusInternalServerError, "Could not retrieve video list")
	}

	data := HomePageData{}
	for _, video := range videoList.Videos {
		data.Videos = append(data.Videos, Video{
			Title:        video.VideoTitle,
			VideoID:      video.VideoID,
			Views:        video.Views,
			AuthorID:     0, // TODO
			AuthorName:   video.AuthorName,
			ThumbnailLoc: video.ThumbnailLoc,
			Rating:       video.Rating,
		})
	}

	return c.Render(http.StatusOK, "home", data)
}
