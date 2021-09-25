package routes

import (
	"context"
	"net/http"
	"strconv"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
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
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	_, err = v.v.ApproveVideo(context.Background(), &videoproto.VideoApproval{VideoID: idInt, UserID: profile.UserID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
