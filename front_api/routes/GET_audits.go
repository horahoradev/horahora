package routes

import (
	"context"
	"net/http"
	"strconv"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /rate/:id where id is video id
// Accepts query parameter "rating" (float)
// Requires authentication
// response: 200 if ok
func (v RouteHandler) handleGetAudits(c echo.Context) error {
	userID := c.Param("id")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		userIDInt = -1 // no user
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	// Only admins can see audit logs
	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	pageNum := getPageNumber(c)

	resp, err := v.u.GetAuditEvents(context.TODO(), &userproto.AuditEventsListRequest{
		Page:   pageNum,
		UserId: userIDInt,
	})
	if err != nil {
		return err
	}

	var data AuditData
	data.Length = len(resp.Events)

	for _, event := range resp.Events {
		data.Events = append(data.Events, AuditEvent{
			ID:        event.Id,
			UserID:    event.User_ID,
			Message:   event.Message,
			Timestamp: event.Timestamp,
		})
	}

	return c.JSON(http.StatusOK, data)
}
