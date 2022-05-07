package routes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/horahoradev/horahora/frontend/internal/config"
	custommiddleware "github.com/horahoradev/horahora/frontend/internal/middleware"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config) {
	r := NewRouteHandler(cfg.VideoClient, cfg.UserClient, cfg.SchedulerClient)

	e.GET("/", r.getHome)
	e.GET("/users/:id", r.getUser)

	e.GET("/videos/:id", r.getVideo)
	e.POST("/rate/:id", r.handleRating)
	e.POST("/approve/:id", r.handleApproval)

	e.GET("/login", getLogin)
	e.POST("/login", r.handleLogin)

	e.GET("/register", getRegister)
	e.POST("/register", r.handleRegister)

	e.GET("/archiverequests", r.getArchiveRequests)
	e.POST("/archiverequests", r.handleArchiveRequest)

	e.GET("/comments/:id", r.getComments)
	e.POST("/comments/", r.handleComment)

	e.POST("/comment_upvotes/", r.handleUpvote)
	e.GET("/upload", getUpload)
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
	L                 LoggedInUserData
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
	Comments          []Comment
	Tags              []string
	RecommendedVideos []Video
	NextVideo         int64
}

type LoggedInUserData struct {
	UserID            int64
	Username          string
	ProfilePictureURL string
}

type ProfileData struct {
	L                 LoggedInUserData
	PaginationData    PaginationData
	UserID            int64
	Username          string
	ProfilePictureURL string
	Videos            []Video
}
type HomePageData struct {
	L              LoggedInUserData
	PaginationData PaginationData
	Videos         []Video
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

func (r RouteHandler) getArchiveRequests(c echo.Context) error {
	data := ArchiveRequestsPageData{}

	addUserProfileInfo(c, &data.L, r.u)

	if data.L.Username == "" {
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
	urlVal := c.FormValue("url")

	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	req := schedulerproto.URLRequest{
		UserID: UserIDInt,
		Url:    urlVal,
	}

	_, err := r.s.DlURL(context.TODO(), &req)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, "/archiverequests")
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
		log.Errorf("Failed to authenticate %s. Err: %s", username, err)
		return c.String(http.StatusForbidden, "Failed to authenticate, invalid credentials")
	}

	return setCookie(c, loginResp.Jwt)
}

func setCookie(c echo.Context, jwt string) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"

	cookie.Value = base64.StdEncoding.EncodeToString([]byte(jwt)) //
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// cookie.SameSite = http.SameSiteStrictMode
	//cookie.Secure = true // set this later

	c.SetCookie(cookie)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (v RouteHandler) getUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	pageNumberInt := getPageNumber(c)

	rank, ok := c.Get(custommiddleware.UserRank).(int32)
	if !ok {
		log.Error("Failed to assert user rank to an int (this should not happen)")
	}
	// doesn't matter if it fails, 0 is a fine default rank
	showUnapproved := false
	if rank > 0 {
		// privileged user, can show unapproved videos
		showUnapproved = true
	}

	videoQueryConfig := videoproto.VideoQueryConfig{
		OrderBy:        videoproto.OrderCategory_upload_date,
		Direction:      videoproto.SortDirection_desc,
		PageNumber:     pageNumberInt,
		SearchVal:      "",
		FromUserID:     idInt,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := v.v.GetVideoList(context.TODO(), &videoQueryConfig)
	if err != nil {
		return err
	}

	getUserReq := userproto.GetUserFromIDRequest{UserID: idInt}

	user, err := v.u.GetUserFromID(context.TODO(), &getUserReq)
	if err != nil {
		return err
	}

	pageRange, err := getPageRange(int(videoList.NumberOfVideos), int(pageNumberInt))
	if err != nil {
		err1 := fmt.Errorf("failed to calculate page range. Err: %s", err)
		log.Error(err1)
		pageRange = []int{1}
	}

	queryStrings := generateQueryParams(pageRange, c)
	data := ProfileData{
		UserID:            idInt,
		Username:          user.Username,
		ProfilePictureURL: "/static/images/placeholder.png",
		PaginationData: PaginationData{
			Pages:                pageRange,
			PathsAndQueryStrings: queryStrings,
			CurrentPage:          int(pageNumberInt),
		},
	}

	for _, video := range videoList.Videos {
		v := Video{
			Title:        video.VideoTitle,
			VideoID:      video.VideoID,
			Views:        video.Views,
			AuthorID:     0, // TODO
			AuthorName:   video.AuthorName,
			ThumbnailLoc: video.ThumbnailLoc,
			Rating:       video.Rating,
		}

		data.Videos = append(data.Videos, v)
	}

	addUserProfileInfo(c, &data.L, v.u)

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

func (v *RouteHandler) getVideo(c echo.Context) error {
	id := c.Param("id")

	// Dumb
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("Could not retrieve video id. Err: %s", err)
	}

	// Increment views first
	viewReq := videoproto.VideoViewing{VideoID: idInt}
	_, err = v.v.ViewVideo(context.Background(), &viewReq)
	if err != nil {
		log.Errorf("Could not increment video views. Err: %s")
	}

	videoReq := videoproto.VideoRequest{
		VideoID: id,
	}

	videoInfo, err := v.v.GetVideo(context.Background(), &videoReq)
	if err != nil {
		log.Errorf("Failed to get video. Err: %s", err)
		return c.String(http.StatusForbidden, "Failed to retrieve video metadata")
	}

	rating := videoInfo.Rating

	// lol
	if math.IsNaN(rating) {
		rating = 0.00
	}

	// TODO: stop copy pasting this...
	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
	}

	recResp, err := v.v.GetVideoRecommendations(context.Background(), &videoproto.RecReq{
		UserId: UserIDInt,
	})
	if err != nil {
		log.Errorf("Could not retrieve recommendations. Err: %s", err)
		// Continue anyway
	}

	var recVideos []Video
	if recResp != nil {
		for _, rec := range recResp.Videos {
			// FIXME: fill other fields after modifying protocol
			vid := Video{
				Title:        rec.VideoTitle,
				VideoID:      rec.VideoID,
				ThumbnailLoc: rec.ThumbnailLoc,
			}

			recVideos = append(recVideos, vid)
		}
	}

	data := VideoDetail{
		L:                 LoggedInUserData{},
		Title:             videoInfo.VideoTitle,
		MPDLoc:            videoInfo.VideoLoc, // FIXME: fix this in videoservice LOL this is embarrassing
		Views:             videoInfo.Views,
		Rating:            rating,
		AuthorID:          videoInfo.AuthorID, // TODO
		Username:          videoInfo.AuthorName,
		UserDescription:   "", // TODO: not implemented yet
		VideoDescription:  videoInfo.Description,
		UserSubscribers:   0, // TODO: not implemented yet
		ProfilePicture:    "/static/images/placeholder.png",
		UploadDate:        videoInfo.UploadDate,
		VideoID:           videoInfo.VideoID,
		Comments:          nil,
		Tags:              videoInfo.Tags,
		RecommendedVideos: recVideos,
	}

	// TODO: awk
	if recVideos != nil && len(recVideos) > 0 {
		data.NextVideo = recVideos[0].VideoID
	}

	addUserProfileInfo(c, &data.L, v.u)

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
		UserID:  UserIDInt,
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

func getPageNumber(c echo.Context) int64 {
	pageNumber := c.QueryParam("page")
	var pageNumberInt int64 = 1

	if pageNumber != "" {
		num, err := strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			log.Errorf("Invalid page number %s, defaulting to 1", pageNumber)
		}
		pageNumberInt = num
	}

	return pageNumberInt
}

func (h *RouteHandler) getHome(c echo.Context) error {
	// SQL injection shouldn't be an issue here, just becomes a list of conditions
	search, err := url.QueryUnescape(c.QueryParam("search"))
	if err != nil {
		return err
	}

	rank, ok := c.Get(custommiddleware.UserRank).(int32)
	if !ok {
		log.Error("Failed to assert user rank to an int (this should not happen)")
	}
	// doesn't matter if it fails, 0 is a fine default rank
	showUnapproved := false
	if rank > 0 {
		// privileged user, can show unapproved videos
		showUnapproved = true
	}

	orderByVal, err := url.QueryUnescape(c.QueryParam("category"))
	if err != nil {
		return err
	}

	// Default
	if orderByVal == "" {
		orderByVal = "upload_date"
	}
	orderBy := videoproto.OrderCategory(videoproto.OrderCategory_value[orderByVal])

	var order videoproto.SortDirection
	orderVal, err := url.QueryUnescape(c.QueryParam("order"))
	if err != nil {
		return err
	}

	if orderVal != "" {
		order = videoproto.SortDirection(videoproto.SortDirection_value[orderVal])
	} else {
		order = videoproto.SortDirection_desc
	}

	pageNumberInt := getPageNumber(c)

	// TODO: if request times out, maybe provide a default list of good videos
	req := videoproto.VideoQueryConfig{
		OrderBy:        orderBy,
		Direction:      order,
		SearchVal:      search,
		PageNumber:     pageNumberInt,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := h.v.GetVideoList(context.TODO(), &req)
	if err != nil {
		log.Errorf("Could not retrieve video list. Err: %s", err)
		return c.String(http.StatusInternalServerError, "Could not retrieve video list")
	}

	pageRange, err := getPageRange(int(videoList.NumberOfVideos), int(pageNumberInt))
	if err != nil {
		err1 := fmt.Errorf("failed to calculate page range. Err: %s", err)
		log.Error(err1)
		pageRange = []int{1}
	}

	queryStrings := generateQueryParams(pageRange, c)

	data := HomePageData{
		PaginationData: PaginationData{
			Pages:                pageRange,
			PathsAndQueryStrings: queryStrings,
			CurrentPage:          int(pageNumberInt),
		},
	}

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

func (v RouteHandler) handleApproval(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	rank, ok := c.Get(custommiddleware.UserRank).(int32)
	if !ok {
		log.Error("Failed to assert user rank to an int (this should not happen)")
	}

	if rank < 1 {
		// privileged user, can show unapproved videos
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	// THERE IS TOO MUCH COPY PASTA HERE!
	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	_, err = v.v.ApproveVideo(context.Background(), &videoproto.VideoApproval{VideoID: idInt, UserID: UserIDInt})
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Video approved")
}

func (r RouteHandler) getComments(c echo.Context) error {
	videoID, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		return err
	}

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return err
	}

	// FIX LATER
	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64 for getComments")
	}

	resp, err := r.v.GetCommentsForVideo(context.Background(), &videoproto.CommentRequest{VideoID: videoIDInt, CurrUserID: UserIDInt})
	if err != nil {
		return err
	}

	commentList := make([]CommentData, 0)

	for _, comment := range resp.Comments {
		commentData := CommentData{
			ID:                 comment.CommentId,
			CreationDate:       comment.CreationDate,
			Content:            comment.Content,
			Username:           comment.AuthorUsername,
			ProfileImage:       comment.AuthorProfileImageUrl,
			VoteScore:          comment.VoteScore,
			CurrUserHasUpvoted: comment.CurrentUserHasUpvoted,
		}
		if comment.ParentId != 0 {
			commentData.ParentID = comment.ParentId
		}

		commentList = append(commentList, commentData)
	}

	return c.JSON(http.StatusOK, &commentList)
}

func (r RouteHandler) handleComment(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return err
	}

	data := c.Request().PostForm
	videoID, err := url.QueryUnescape(data.Get("video_id"))
	if err != nil {
		return err
	}

	userID, err := url.QueryUnescape(data.Get("user_id"))
	if err != nil {
		return err
	}

	content, err := url.QueryUnescape(data.Get("content"))
	if err != nil {
		return err
	}

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return err
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	parentIDInt, _ := getAsInt64(data, "parent")
	//if err != nil {
	//	// nothing
	//}

	_, err = r.v.MakeComment(context.Background(), &videoproto.VideoComment{
		UserId:        userIDInt,
		VideoId:       videoIDInt,
		Comment:       content,
		ParentComment: parentIDInt,
	})

	if err != nil {
		return err
	}

	return nil
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

func (r RouteHandler) handleUpvote(c echo.Context) error {
	// DUMB!
	err := c.Request().ParseForm()
	if err != nil {
		return err
	}

	data := c.Request().PostForm

	commentID, err := getAsInt64(data, "comment_id")
	if err != nil {
		return err
	}

	userID, err := getAsInt64(data, "user_id")
	if err != nil {
		return err
	}

	hasUpvoted, err := getAsBool(data, "user_has_upvoted")
	if err != nil {
		return err
	}

	_, err = r.v.MakeCommentUpvote(context.Background(), &videoproto.CommentUpvote{
		CommentId: commentID,
		UserId:    userID,
		IsUpvote:  hasUpvoted,
	})

	return err
}

const (
	videoKey               = "file[0]"
	thumbnailKey           = "file[1]"
	MINIMUM_NUMBER_OF_TAGS = 5
	fileUploadChunkSize    = 1024 * 1024
)

func (r RouteHandler) upload(c echo.Context) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return err
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	tagsInt := c.FormValue("tags")

	var tags []string
	err = json.Unmarshal([]byte(tagsInt), &tags)
	if err != nil {
		return err
	}

	thumbFileHeader, err := c.FormFile(thumbnailKey)
	if err != nil {
		return err
	}

	thumbFile, err := thumbFileHeader.Open()
	if err != nil {
		return err
	}

	videoFileHeader, err := c.FormFile(videoKey)
	if err != nil {
		return err
	}

	videoFile, err := videoFileHeader.Open()
	if err != nil {
		return err
	}

	// TODO: rewrite me, this isn't memory efficient
	videoBytes, err := ioutil.ReadAll(videoFile)
	if err != nil {
		return err
	}

	thumbBytes, err := ioutil.ReadAll(thumbFile)
	if err != nil {
		return err
	}

	log.Infof("Title: %s", title)
	log.Infof("Description: %s", description)
	log.Infof("Tags: %s", tags)
	log.Infof("Video len: %d thumb len: %d", len(videoBytes), len(thumbBytes))

	uploadClient, err := r.v.UploadVideo(context.Background())
	if err != nil {
		return err
	}

	metaChunk := &videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Meta{
			Meta: &videoproto.InputFileMetadata{
				Title:       title,
				Description: description,
				//AuthorUID:            "",
				//OriginalVideoLink:    "",
				//AuthorUsername:       "",
				//OriginalSite:         0,
				//OriginalID:           "",
				DomesticAuthorID: userID,
				Tags:             tags,
				Thumbnail:        thumbBytes,
			},
		},
	}

	err = uploadClient.Send(metaChunk)
	if err != nil {
		return err
	}

	for byteInd := 0; byteInd < len(videoBytes); byteInd += fileUploadChunkSize {
		videoByteSlice := videoBytes[byteInd:min(len(videoBytes), byteInd+fileUploadChunkSize)]
		log.Infof("uploading byte %d", byteInd)
		videoChunk := &videoproto.InputVideoChunk{
			Payload: &videoproto.InputVideoChunk_Content{
				Content: &videoproto.FileContent{
					Data: videoByteSlice,
				},
			},
		}

		err = uploadClient.Send(videoChunk)
		if err != nil {
			return err
		}
	}

	resp, err := uploadClient.CloseAndRecv()
	if err != nil {
		return err
	}

	// Redirect to the new video
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/videos/%d", resp.VideoID))
}

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

func getUpload(c echo.Context) error {
	return c.Render(http.StatusOK, "upload", nil)
}
