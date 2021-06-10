package routes

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	custommiddleware "github.com/SEAPUNK/horahora/front_api/middleware"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (v RouteHandler) handleApproval(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	rank, ok := c.Get(custommiddleware.UserRankKey).(int32)
	if !ok {
		log.Error("Failed to assert user rank to an int (this should not happen)")
	}

	if rank < 1 {
		// privileged user, can show unapproved videos
		// TODO(ivan): status forbidden
		return errors.New("Insufficient user status")
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

	return c.JSON(http.StatusOK, nil)
}
