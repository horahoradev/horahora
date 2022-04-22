package routes

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/horahoradev/horahora/front_api/config"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
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
	e.GET("/currentuserprofile/", r.getCurrentUserProfile)

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
	Title             string
	MPDLoc            string
	Views             uint64
	Rating            float64
	VideoID           int64
	AuthorID          int64
	Username          string
	UserDescription   string
	VideoDescription  string
	UserSubscribers   uint64
	ProfilePicture    string
	UploadDate        string // should be a datetime
	Tags              []string
	RecommendedVideos []Video
	L                 *LoggedInUserData
}

type ProfileData struct {
	PaginationData    PaginationData
	UserID            int64
	Username          string
	ProfilePictureURL string
	Videos            []Video
	L                 *LoggedInUserData
}

type PaginationData struct {
	NumberOfItems int
	CurrentPage   int
}

type ArchiveRequestsPageData struct {
	ArchivalRequests []*schedulerproto.ContentArchivalEntry
	ArchivalEvents   []*schedulerproto.ArchivalEvent
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
