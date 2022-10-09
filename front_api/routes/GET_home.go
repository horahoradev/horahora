package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HomePageData struct {
	PaginationData PaginationData
	Videos         []Video
}

// route: GET /home[?seach=val][&page=x]
// The query string val search is accepted, and will return videos whose title, tags, or description contains the search term. Inclusion and exclusion is supported, e.g. include1 include2 -exclude1
// A query string val for the page number, starting at 1, is also accepted.
// Response is of the form: {"PaginationData":{"PathsAndQueryStrings":["/home?page=1"],"Pages":[1],"CurrentPage":1},"Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":6,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}
// For pagination data, the fields pages and PathsAndQueryStrings will always have the same length, and have corresponding values
// Response is of the form: {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":9,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}
func (h *RouteHandler) getHome(c echo.Context) error {
	// TODO: verify no sql injection lol
	search, err := url.QueryUnescape(c.QueryParam("search"))
	if err != nil {
		return err
	}

	profileInfo, err := h.getUserProfileInfo(c)
	if err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	// doesn't matter if it fails, 0 is a fine default rank
	showUnapproved := false
	if profileInfo != nil && profileInfo.Rank > 0 {
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

	orderCat, ok := videoproto.OrderCategory_value[orderByVal]
	if !ok {
		return fmt.Errorf("invalid category supplied: %s", orderByVal)
	}

	orderBy := videoproto.OrderCategory(orderCat)

	orderVal, err := url.QueryUnescape(c.QueryParam("order"))
	if err != nil {
		return err
	}
	var order videoproto.SortDirection
	if orderVal == "" {
		order = videoproto.SortDirection_desc // Default to desc
	} else {
		order = videoproto.SortDirection(videoproto.SortDirection_value[orderVal])
	}

	pageNumber := getPageNumber(c)

	// TODO: if request times out, maybe provide a default list of good videos
	req := videoproto.VideoQueryConfig{
		OrderBy:        orderBy,
		Direction:      order,
		SearchVal:      search,
		PageNumber:     pageNumber,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := h.v.GetVideoList(context.TODO(), &req)
	if err != nil {
		log.Errorf("Could not retrieve video list. Err: %s", err)
		return errors.New("Could not retrieve video list")
	}

	data := HomePageData{
		PaginationData: PaginationData{
			NumberOfItems: int(videoList.NumberOfVideos),
			CurrentPage:   int(pageNumber),
		},
	}

	data.Videos = []Video{}
	for _, video := range videoList.Videos {
		data.Videos = append(data.Videos, Video{
			Title:        video.VideoTitle,
			VideoID:      video.VideoID,
			Views:        video.Views,
			AuthorID:     video.AuthorID, // TODO
			AuthorName:   video.AuthorName,
			ThumbnailLoc: video.ThumbnailLoc,
			Rating:       video.Rating,
		})
	}

	return c.JSON(http.StatusOK, data)
}
