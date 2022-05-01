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

// Route: POST /approve/:id
// Requires authentication
// Allows the user, if sufficiently high rank, to approve of a video and allow it to be shown to regular users.
// Response: 200 if okay
func (v RouteHandler) handleDelete(c echo.Context) error {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 64) // just make sure we can parse it
	if err != nil {
		return err
	}

	// Requesting profile
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	// Make an audit event even if they don't pass the permission check
	_, err = v.u.AddAuditEvent(context.TODO(), &userproto.NewAuditEventRequest{
		Message: fmt.Sprintf("User attempted to delete video id %s", id),
		User_ID: profile.UserID,
	})
	if err != nil {
		return err // If the audit event can't be created, fail the operation
	}

	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	_, err = v.v.DeleteVideo(context.Background(), &videoproto.VideoDeletionReq{VideoID: id})

	return err
}
