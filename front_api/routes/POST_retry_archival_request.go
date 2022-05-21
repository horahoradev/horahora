package routes

import (
	"context"
	"fmt"
	"strconv"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

func (v RouteHandler) handleRetryArchivalRequest(c echo.Context) error {
	downloadID := c.FormValue("download_id")
	downloadIDInt, err := strconv.ParseInt(downloadID, 10, 64) // just make sure we can parse it
	if err != nil {
		return err
	}

	// Requesting profile
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	// We want an audit log entry for this one
	_, err = v.u.AddAuditEvent(context.TODO(), &userproto.NewAuditEventRequest{
		Message: fmt.Sprintf("User attempted to retry archival request %d", downloadIDInt),
		User_ID: profile.UserID,
	})
	if err != nil {
		return err // If the audit event can't be created, fail the operation
	}

	_, err = v.s.RetryArchivalRequestDownloadss(context.TODO(), &schedulerproto.RetryRequest{UserID: uint64(profile.UserID),
		DownloadID: uint64(downloadIDInt)})

	return err
}
