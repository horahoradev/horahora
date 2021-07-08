package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	custommiddleware "github.com/SEAPUNK/horahora/front_api/middleware"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HomePageData struct {
	L              LoggedInUserData
	PaginationData PaginationData
	Videos         []Video
}

func (h *RouteHandler) getHome(c echo.Context) error {
	// TODO: verify no sql injection lol
	tag, err := url.QueryUnescape(c.QueryParam("tag"))
	if err != nil {
		return err
	}

	rank, ok := c.Get(custommiddleware.UserRankKey).(int32)
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

	orderVal, err := url.QueryUnescape(c.QueryParam("order"))
	if err != nil {
		return err
	}
	order := videoproto.SortDirection(videoproto.SortDirection_value[orderVal])

	pageNumber := c.QueryParam("page")
	var pageNumberInt int64 = 1

	if pageNumber != "" {
		num, err := strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			log.Errorf("Invalid page number %s, defaulting to 1", pageNumber)
		}
		pageNumberInt = num
	}

	// TODO: if request times out, maybe provide a default list of good videos
	req := videoproto.VideoQueryConfig{
		OrderBy:        orderBy,
		Direction:      order,
		SearchVal:      tag,
		PageNumber:     pageNumberInt,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := h.v.GetVideoList(context.TODO(), &req)
	if err != nil {
		log.Errorf("Could not retrieve video list. Err: %s", err)
		return errors.New("Could not retrieve video list")
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

	return c.JSON(http.StatusOK, data)
}
