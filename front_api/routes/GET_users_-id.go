package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

// route: GET /users/:id where id is the user id
// A query string val for the page number, starting at 1, is also accepted.
// Response is of the form:{"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}
// For pagination data, the fields pages and PathsAndQueryStrings will always have the same length, and have corresponding values
func (v RouteHandler) getUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	getUserReq := userproto.GetUserFromIDRequest{UserID: idInt}

	profile, err := v.u.GetUserFromID(context.TODO(), &getUserReq)
	if err != nil {
		return fmt.Errorf("Get user from ID: %s", err)
	}

	userProfile, err := v.getUserProfileInfo(c)
	if err != nil {
		return fmt.Errorf("Get user profile info: %s", err)
	}

	// doesn't matter if it fails, 0 is a fine default rank
	showUnapproved := false
	if userProfile.Rank > 0 {
		// privileged user, can show unapproved videos
		showUnapproved = true
	}

	pageNumber := getPageNumber(c)

	videoQueryConfig := videoproto.VideoQueryConfig{
		OrderBy:        videoproto.OrderCategory_upload_date,
		Direction:      videoproto.SortDirection_desc,
		PageNumber:     pageNumber,
		SearchVal:      "",
		FromUserID:     idInt,
		ShowUnapproved: showUnapproved,
	}

	videoList, err := v.v.GetVideoList(context.TODO(), &videoQueryConfig)
	if err != nil {
		return fmt.Errorf("Get video list: %s", err)
	}

	// TODO: 0 results in all videos, fix for admin user?
	data := ProfileData{
		UserID:            idInt,
		Username:          profile.Username,
		L:                 userProfile,
		ProfilePictureURL: "/static/images/placeholder1.jpg",
		PaginationData: PaginationData{
			NumberOfItems: int(videoList.NumberOfVideos),
			CurrentPage:   int(pageNumber),
		},
	}

	data.Videos = []Video{}
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
