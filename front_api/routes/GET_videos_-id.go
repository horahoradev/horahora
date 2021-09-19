package routes

import (
	"context"
	"math"
	"net/http"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

func (v *RouteHandler) getVideo(c echo.Context) error {
	id := c.Param("id")

	// Dumb
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// Increment views first
	viewReq := videoproto.VideoViewing{VideoID: idInt}
	_, err = v.v.ViewVideo(context.Background(), &viewReq)
	if err != nil {
		return err
	}

	videoReq := videoproto.VideoRequest{
		VideoID: id,
	}

	videoInfo, err := v.v.GetVideo(context.Background(), &videoReq)
	if err != nil {
		return err
	}

	rating := videoInfo.Rating

	// lol
	if math.IsNaN(rating) {
		rating = 0.00
	}

	data := VideoDetail{
		L:                LoggedInUserData{},
		Title:            videoInfo.VideoTitle,
		MPDLoc:           videoInfo.VideoLoc, // FIXME: fix this in videoservice LOL this is embarrassing
		Views:            videoInfo.Views,
		Rating:           rating,
		AuthorID:         videoInfo.AuthorID, // TODO
		Username:         videoInfo.AuthorName,
		UserDescription:  "", // TODO: not implemented yet
		VideoDescription: videoInfo.Description,
		UserSubscribers:  0, // TODO: not implemented yet
		ProfilePicture:   "/static/images/placeholder1.jpg",
		UploadDate:       videoInfo.UploadDate,
		VideoID:          videoInfo.VideoID,
		Comments:         nil,
		Tags:             videoInfo.Tags,
	}

	addUserProfileInfo(c, &data.L, v.u)

	return c.JSON(http.StatusOK, data)
}
