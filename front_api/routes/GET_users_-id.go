package routes

import (
	"context"
	"fmt"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func (v RouteHandler) getUser(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	// doesn't matter if it fails, 0 is a fine default rank
	showUnapproved := false
	if profile.Rank > 0 {
		// privileged user, can show unapproved videos
		showUnapproved = true
	}

	pageNumber := getPageNumber(c)

	videoQueryConfig := videoproto.VideoQueryConfig{
		OrderBy:        videoproto.OrderCategory_upload_date,
		Direction:      videoproto.SortDirection_desc,
		PageNumber:     pageNumber,
		SearchVal:      "",
		FromUserID:     profile.UserID,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := v.v.GetVideoList(context.TODO(), &videoQueryConfig)
	if err != nil {
		return err
	}

	pageRange, err := getPageRange(int(videoList.NumberOfVideos), int(pageNumber))
	if err != nil {
		err1 := fmt.Errorf("failed to calculate page range. Err: %s", err)
		log.Error(err1)
		pageRange = []int{1}
	}

	queryStrings := generateQueryParams(pageRange, c)
	data := ProfileData{
		UserID:            profile.UserID,
		Username:          profile.Username,
		ProfilePictureURL: "/static/images/placeholder1.jpg",
		PaginationData: PaginationData{
			Pages:                pageRange,
			PathsAndQueryStrings: queryStrings,
			CurrentPage:          int(pageNumber),
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
	
	return c.JSON(http.StatusOK, data)
}
