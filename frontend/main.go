package main

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/horahoradev/horahora/frontend/internal/config"
	custommiddleware "github.com/horahoradev/horahora/frontend/internal/middleware"
	"github.com/horahoradev/horahora/frontend/internal/templates"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	r := NewRouteHandler(cfg.VideoClient, cfg.UserClient, cfg.SchedulerClient)

	e.GET("/", r.getHome)
	e.GET("/users/:id", getUser)

	e.GET("/videos/:id", r.getVideo)
	e.POST("/rate/:id", r.handleRating)

	e.GET("/login", getLogin)
	e.POST("/login", r.handleLogin)

	e.GET("/register", getRegister)
	e.POST("/register", r.handleRegister)

	e.GET("/archiverequests", r.getArchiveRequests)
	e.POST("/archiverequests", r.handleArchiveRequest)

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
	VideoID         int64
	AuthorID        int64
	Username        string
	UserDescription string
	UserSubscribers uint64
	ProfilePicture  string
	UploadDate      string // should be a datetime
	Comments        []Comment
	Tags            []string
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

type ArchiveRequestsPageData struct {
	L                LoggedInUserData
	ArchivalRequests []*schedulerproto.ContentArchivalEntry
}

func (r RouteHandler) getArchiveRequests(c echo.Context) error {
	data := ArchiveRequestsPageData{}

	addUserProfileInfo(c, &data.L, r.u)

	if data.L.UserID == 0 {
		// User isn't logged in
		// TODO: move this to a middleware somehow
		return c.String(http.StatusForbidden, "Must be logged in")
	}

	resp, err := r.s.ListArchivalEntries(context.TODO(), &schedulerproto.ListArchivalEntriesRequest{UserID: data.L.UserID})
	if err != nil {
		return err
	}

	data.ArchivalRequests = resp.Entries

	return c.Render(http.StatusOK, "archiveRequests", data)
}

func (r RouteHandler) handleArchiveRequest(c echo.Context) error {
	website := c.FormValue("website")
	contentType := c.FormValue("contentType")
	contentValue := c.FormValue("contentValue")

	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	websiteEnumVal, ok := schedulerproto.SupportedSite_value[website]
	if !ok {
		return errors.New("site not found")
	}

	supportedWebsite := schedulerproto.SupportedSite(websiteEnumVal)

	// FIXME: this is dumb. Fix this to use schedulerproto consts after switching to string instead of enum
	switch contentType {
	case "tag":
		req := schedulerproto.TagRequest{
			UserID:   UserIDInt,
			Website:  supportedWebsite, // FIXME: placeholder, see above
			TagValue: contentValue,
		}

		_, err := r.s.DlTag(context.TODO(), &req)
		if err != nil {
			return err
		}
	case "channel":
		req := schedulerproto.ChannelRequest{
			Website:   supportedWebsite,
			ChannelID: contentValue,
		}

		_, err := r.s.DlChannel(context.TODO(), &req)
		if err != nil {
			return err
		}
	case "playlist":
		req := schedulerproto.PlaylistRequest{
			Website:    supportedWebsite,
			PlaylistID: contentValue,
		}
		_, err := r.s.DlPlaylist(context.TODO(), &req)
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid content type")
	}

	return c.String(http.StatusOK, "Archive request submitted successfully")
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func getRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func (r RouteHandler) handleRegister(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	registrationReq := userproto.RegisterRequest{
		Password: password,
		Username: username,
		Email:    email,
	}

	regisResp, err := r.u.Register(context.Background(), &registrationReq)
	if err != nil {
		return err
	}

	// TODO: use registration JWT to auth

	return setCookie(c, regisResp.Jwt)
}

func (r RouteHandler) handleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO: grpc auth goes here
	loginReq := &userproto.LoginRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := r.u.Login(context.Background(), loginReq)
	if err != nil {
		return err
	}

	return setCookie(c, loginResp.Jwt)
}

func setCookie(c echo.Context, jwt string) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"

	cookie.Value = base64.StdEncoding.EncodeToString([]byte(jwt)) //
	cookie.Expires = time.Now().Add(24 * time.Hour)

	cookie.SameSite = http.SameSiteStrictMode
	//cookie.Secure = true // set this later

	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Login successful.")
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

	// addUserProfileInfo(c, &data.L, )
	// TODO: get the rest of data

	return c.Render(http.StatusOK, "profile", data)
}

type RouteHandler struct {
	v videoproto.VideoServiceClient
	u userproto.UserServiceClient
	s schedulerproto.SchedulerClient
}

func NewRouteHandler(v videoproto.VideoServiceClient, u userproto.UserServiceClient, s schedulerproto.SchedulerClient) *RouteHandler {
	return &RouteHandler{
		v: v,
		u: u,
		s: s,
	}
}

func (v *RouteHandler) getVideo(c echo.Context) error {
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

	rating := videoInfo.Rating

	// lol
	if math.IsNaN(rating) {
		rating = 0.00
	}

	data := VideoDetail{
		L:               LoggedInUserData{},
		Title:           videoInfo.VideoTitle,
		MPDLoc:          videoInfo.VideoLoc, // FIXME: fix this in videoservice LOL this is embarrassing
		Views:           videoInfo.Views,
		Rating:          rating,
		AuthorID:        videoInfo.AuthorID, // TODO
		Username:        videoInfo.AuthorName,
		UserDescription: "", // TODO: not implemented yet
		UserSubscribers: 0,  // TODO: not implemented yet
		ProfilePicture:  "/static/images/placeholder1.jpg",
		UploadDate:      videoInfo.UploadDate,
		VideoID:         videoInfo.VideoID,
		Comments:        nil,
		Tags: videoInfo.
	}

	addUserProfileInfo(c, &data.L, v.u)

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

func (v RouteHandler) handleRating(c echo.Context) error {
	videoID := c.Param("id")
	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Error("Could not assert videoID to int64")
		return errors.New("could not assert videoID to int64")
	}

	ratings, ok := c.QueryParams()["rating"]
	if !ok {
		return errors.New("no rating in query string")
	}

	rating, err := strconv.ParseFloat(ratings[0], 64)
	if err != nil {
		return err
	}

	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	rateReq := videoproto.VideoRating{
		UserID:  string(UserIDInt),
		VideoID: videoIDInt,
		Rating:  float32(rating),
	}

	_, err = v.v.RateVideo(context.TODO(), &rateReq)
	return err
}

type HomeHandler struct {
	videoClient videoproto.VideoServiceClient
	userClient  userproto.UserServiceClient
}

func (h *RouteHandler) getHome(c echo.Context) error {
	// TODO: if request times out, maybe provide a default list of good videos
	req := videoproto.VideoQueryConfig{
		OrderBy:    videoproto.OrderCategory_upload_date,
		Direction:  videoproto.SortDirection_desc,
		PageNumber: 0,
	}

	videoList, err := h.v.GetVideoList(context.TODO(), &req)
	if err != nil {
		log.Errorf("Could not retrieve video list. Err: %s", err)
		return c.String(http.StatusInternalServerError, "Could not retrieve video list")
	}

	data := HomePageData{}

	addUserProfileInfo(c, &data.L, h.u)
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

func addUserProfileInfo(c echo.Context, l *LoggedInUserData, client userproto.UserServiceClient) {
	id := c.Get(custommiddleware.UserIDKey)

	idInt, ok := id.(int64)
	if !ok {
		log.Error("Could not assert id to int64")
		return
	}

	if idInt == 0 {
		return // User isn't logged in
	}

	getUserReq := userproto.GetUserFromIDRequest{
		UserID: idInt,
	}

	userResp, err := client.GetUserFromID(context.TODO(), &getUserReq)
	if err != nil {
		log.Error(err)
		return
	}

	l.Username = userResp.Username
	// l.ProfilePictureURL = userResp. // TODO
	l.UserID = idInt
}
