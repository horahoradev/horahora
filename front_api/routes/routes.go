package routes

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/horahoradev/horahora/front_api/config"
	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

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

func SetupRoutes(e *echo.Echo, cfg *config.Config) {
	r := NewRouteHandler(cfg.VideoClient, cfg.UserClient, cfg.SchedulerClient)

	e.GET("/home", r.getHome)
	e.GET("/users/:id", r.getUser)

	e.GET("/tag/:tag", r.getTag)

	e.GET("/videos/:id", r.getVideo)
	e.POST("/rate/:id", r.handleRating)
	e.POST("/approve/:id", r.handleApproval)

	e.POST("/login", r.handleLogin)
	e.POST("/register", r.handleRegister)
	e.POST("/logout", r.handleLogout)

	e.GET("/archiverequests", r.getArchiveRequests)
	e.POST("/archiverequests", r.handleArchiveRequest)

	e.GET("/comments/:id", r.getComments)
	e.POST("/comments/", r.handleComment)

	e.POST("/comment_upvotes/", r.handleUpvote)
	e.POST("/upload", r.upload)
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
	L                LoggedInUserData
	Title            string
	MPDLoc           string
	Views            uint64
	Rating           float64
	VideoID          int64
	AuthorID         int64
	Username         string
	UserDescription  string
	VideoDescription string
	UserSubscribers  uint64
	ProfilePicture   string
	UploadDate       string // should be a datetime
	Comments         []Comment
	Tags             []string
}

type LoggedInUserData struct {
	UserID            int64
	Username          string
	ProfilePictureURL string
	Rank              int32
}

type ProfileData struct {
	L                 LoggedInUserData
	PaginationData    PaginationData
	UserID            int64
	Username          string
	ProfilePictureURL string
	Videos            []Video
}

type PaginationData struct {
	PathsAndQueryStrings []string
	Pages                []int
	CurrentPage          int
}

type ArchiveRequestsPageData struct {
	L                LoggedInUserData
	ArchivalRequests []*schedulerproto.ContentArchivalEntry
}

func setCookie(c echo.Context, jwt string) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"

	cookie.Value = base64.StdEncoding.EncodeToString([]byte(jwt)) //
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// cookie.SameSite = http.SameSiteStrictMode
	//cookie.Secure = true // set this later

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, nil)
}

func generateQueryParams(pageRange []int, c echo.Context) []string {
	// This is ugly
	var queryStrings []string
	for _, page := range pageRange {
		p := strconv.FormatInt(int64(page), 10)

		c.QueryParams().Set("page", p) // FIXME
		queryStrings = append(queryStrings, c.Request().URL.Path+"?"+c.QueryParams().Encode())
	}

	return queryStrings
}

func getCurrentUserID(c echo.Context) (int64, error) {
	id := c.Get(custommiddleware.UserIDKey)

	idInt, ok := id.(int64)
	if !ok {
		log.Error("Could not assert id to int64")
		return 0, errors.New("could not assert id to int64")
	}

	return idInt, nil
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
	l.Rank = int32(userResp.Rank)
}

type CommentData struct {
	ID                 int64  `json:"id"`
	CreationDate       string `json:"created"`
	Content            string `json:"content"`
	Username           string `json:"fullname"`
	ProfileImage       string `json:"profile_picture_url"`
	VoteScore          int64  `json:"upvote_count"`
	CurrUserHasUpvoted bool   `json:"user_has_upvoted"`
	ParentID           int64  `json:"parent,omitempty"`
}

const (
	videoKey               = "file[0]"
	thumbnailKey           = "file[1]"
	MINIMUM_NUMBER_OF_TAGS = 5
	fileUploadChunkSize    = 1024 * 1024
)

func getAsInt64(data url.Values, key string) (int64, error) {
	val, err := url.QueryUnescape(data.Get(key))
	if err != nil {
		return 0, err
	}

	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}

func getAsBool(data url.Values, key string) (bool, error) {
	val, err := url.QueryUnescape(data.Get(key))
	if err != nil {
		return false, err
	}

	valBool, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	return valBool, nil
}
