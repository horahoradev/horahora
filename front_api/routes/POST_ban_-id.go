package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /approve/:id
// Requires authentication
// Allows the user, if sufficiently high rank, to approve of a video and allow it to be shown to regular users.
// Response: 200 if okay
func (v RouteHandler) handleBan(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// Requesting profile
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	req := userproto.BanUserRequest{
		UserID: idInt,
	}

	_, err = v.u.BanUser(context.Background(), &req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("could not ban user: %s", err))
	}

	return c.JSON(http.StatusOK, nil)
}
