package routes

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Route: POST /rate/:id where id is video id
// Accepts query parameter "rating" (float)
// Requires authentication
// response: 200 if ok
func (v RouteHandler) handleRating(c echo.Context) error {
	videoID := c.Param("id")
	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Error("Could not assert videoID to int64")
		return errors.New("could not assert videoID to int64")
	}

	ratingStr := c.FormValue("rating")

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return err
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	log.Info(videoIDInt)

	rateReq := videoproto.VideoRating{
		UserID:  profile.UserID,
		VideoID: videoIDInt,
		Rating:  float32(rating),
	}

	_, err = v.v.RateVideo(context.TODO(), &rateReq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
