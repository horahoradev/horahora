package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

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

	// Make an audit event even if they don't pass the permission check
	_, err = v.u.AddAuditEvent(context.TODO(), &userproto.NewAuditEventRequest{
		Message: fmt.Sprintf("User attempted to ban user id %d", idInt),
		User_ID: profile.UserID,
	})
	if err != nil {
		return err // If the audit event can't be created, fail the operation
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
