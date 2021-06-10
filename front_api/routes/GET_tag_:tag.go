package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	custommiddleware "github.com/SEAPUNK/horahora/front_api/middleware"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (v RouteHandler) getTag(c echo.Context) error {
	tag, err := url.QueryUnescape(c.Param("tag"))
	if err != nil {
		return err
	}

	pageNumber := c.QueryParam("page")
	var pageNumberInt int64 = 1

	if pageNumber != "" {
		num, err := strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			log.Errorf("Invalid page number %s, defaulting to 1", pageNumber)
		}
		pageNumberInt = num
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

	videoQueryConfig := videoproto.VideoQueryConfig{
		OrderBy:        videoproto.OrderCategory_upload_date,
		Direction:      videoproto.SortDirection_desc,
		PageNumber:     pageNumberInt,
		ContainsTag:    tag,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := v.v.GetVideoList(context.TODO(), &videoQueryConfig)
	if err != nil {
		return err
	}

	pageRange, err := getPageRange(int(videoList.NumberOfVideos), int(pageNumberInt))
	if err != nil {
		err1 := fmt.Errorf("failed to calculate page range. Err: %s", err)
		log.Error(err1)
		pageRange = []int{1}
	}

	// TODO: copy pasta is very bad

	queryStrings := generateQueryParams(pageRange, c)

	data := HomePageData{
		PaginationData: PaginationData{
			Pages:                pageRange,
			PathsAndQueryStrings: queryStrings,
			CurrentPage:          int(pageNumberInt),
		},
	}

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

	addUserProfileInfo(c, &data.L, v.u)
	return c.JSON(http.StatusOK, data)
}
