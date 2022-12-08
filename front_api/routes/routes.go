package routes

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/horahoradev/horahora/front_api/config"
	partyproto "github.com/horahoradev/horahora/partyservice/protocol"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"

	"github.com/labstack/echo/v4"
)

type RouteHandler struct {
	v   videoproto.VideoServiceClient
	u   userproto.UserServiceClient
	s   schedulerproto.SchedulerClient
	p   partyproto.PartyserviceClient
	srv *socketio.Server
}

func NewRouteHandler(v videoproto.VideoServiceClient, u userproto.UserServiceClient, s schedulerproto.SchedulerClient, p partyproto.PartyserviceClient, srv *socketio.Server) *RouteHandler {
	return &RouteHandler{
		v:   v,
		u:   u,
		s:   s,
		p:   p,
		srv: srv,
	}
}

func SetupRoutes(e *echo.Echo, cfg *config.Config, srv *socketio.Server) {
	r := NewRouteHandler(cfg.VideoClient, cfg.UserClient, cfg.SchedulerClient, cfg.PartyClient, srv)

	e.GET("/api/home", r.getHome)
	e.GET("/api/users/:id", r.getUser)
	e.GET("/api/currentuserprofile/", r.getCurrentUserProfile)
	e.GET("/api/auditevents/:id", r.handleGetAudits)
	e.GET("/api/downloadsinprogress", r.handleGetDownloadsInProgress)

	e.GET("/api/videos/:id", r.getVideo)
	e.POST("/api/rate/:id", r.handleRating)
	e.POST("/api/approve/:id", r.handleApproval)

	e.POST("/api/login", r.handleLogin)
	e.POST("/api/revoltlogin", r.handleRevoltLogin)
	e.POST("/api/register", r.handleRegister)
	e.POST("/api/logout", r.handleLogout)

	e.GET("/api/archiverequests", r.getArchiveRequests)
	e.GET("/api/archiveevents/:id", r.getArchiveEvents)
	e.POST("/api/archiverequests", r.handleArchiveRequest)
	e.POST("/api/delete-archiverequest", r.handleDeleteArchivalRequest)
	e.POST("/api/retry-archiverequest", r.handleRetryArchivalRequest)

	e.GET("/api/comments/:id", r.getComments)
	e.POST("/api/comments/", r.handleComment)

	e.POST("/api/comment_upvotes/", r.handleUpvote)
	e.POST("/api/upload", r.upload)

	e.POST("/api/ban/:id", r.handleBan)
	e.POST("/api/delete/:id", r.handleDelete)
	e.POST("/api/setrank/:userid/:rank", r.handleSetRank)
	e.POST("/api/password-reset", r.handlePasswordReset)

	// Watch party stuff goes here
	e.POST("/api/newwatchparty/:id", r.handleNewWatchParty)
	e.POST("/api/joinwatchparty/:id", r.handleJoinWatchParty)
	e.POST("/api/heartbeat", r.handleHeartbeat)
	e.GET("/api/partystate/:id", r.handleGetPartyState)
	e.POST("/api/addvideo/:id", r.handleAddVideo)
	e.POST("/api/nextvideo/:id", r.handleNextVideo)

	e.GET("/api/unapprovedvideos", r.getUnapprovedVideos)
	e.POST("/api/approvevideosvideo/:id", r.handleSchedulerVideoApproavl)
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

type VideoInProgress struct {
	Website  string
	VideoID  string
	DlStatus string
}

type UnapprovedVideo struct {
	URL     string
	VideoID string
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
	Banned            bool
	L                 *LoggedInUserData
}

type AuditEvent struct {
	ID        int64
	UserID    int64
	Message   string
	Timestamp string
}

type AuditData struct {
	Length int
	Events []AuditEvent
}

type PaginationData struct {
	NumberOfItems int
	CurrentPage   int
}

type ArchiveRequestsPageData struct {
	ArchivalRequests []ArchivalRequest
}

type ArchiveEventsData struct {
	ArchivalEvents []*schedulerproto.ArchivalEvent
}

type ArchivalRequest struct {
	UserID               int64
	Url                  string
	ArchivedVideos       uint64
	CurrentTotalVideos   uint64
	LastSynced           string
	BackoffFactor        uint32
	DownloadID           uint64
	UndownloadableVideos uint64
}

func setCookie(c echo.Context, jwt string) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"

	cookie.Value = base64.StdEncoding.EncodeToString([]byte(jwt)) //
	cookie.Expires = time.Now().Add(24 * time.Hour)

	cookie.SameSite = http.SameSiteStrictMode
	//cookie.Secure = true // set this later
	cookie.HttpOnly = false

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
