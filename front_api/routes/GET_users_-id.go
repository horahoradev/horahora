package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (v RouteHandler) getUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// TODO: reduce copy pasta
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
		ProfilePictureURL: "/static/images/placeholder1.jpg",
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

	return c.JSON(http.StatusOK, data)
}
