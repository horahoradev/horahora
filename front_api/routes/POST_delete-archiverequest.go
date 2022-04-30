package routes

import (
	"context"
	"net/http"
	"strconv"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"

	"github.com/labstack/echo/v4"
)

func (v RouteHandler) handleDeleteArchivalRequest(c echo.Context) error {
	// Requesting profile
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank < 1 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	downloadID := c.FormValue("download_id")
	downloadIDInt, err := strconv.ParseInt(downloadID, 10, 64) // just make sure we can parse it
	if err != nil {
		return err
	}

	_, err = v.s.DeleteArchivalRequest(context.TODO(), &schedulerproto.DeletionRequest{UserID: uint64(profile.UserID), DownloadID: uint64(downloadIDInt)})

	return err
}
