package routes

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

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
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
