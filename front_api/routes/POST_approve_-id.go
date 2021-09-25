package routes

import (
	"context"
	"net/http"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)


// Route: POST /approve/:id
// Requires authentication
// Allows the user, if sufficiently high rank, to approve of a video and allow it to be shown to regular users.
// Response: 200 if okay
func (v RouteHandler) handleApproval(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank < 1 {
		// privileged user, can show unapproved videos
		// TODO(ivan): status forbidden
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	_, err = v.v.ApproveVideo(context.Background(), &videoproto.VideoApproval{VideoID: idInt, UserID: profile.UserID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
