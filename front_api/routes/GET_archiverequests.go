package routes

import (
	"context"
	"errors"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) getArchiveRequests(c echo.Context) error {
	data := ArchiveRequestsPageData{}

	addUserProfileInfo(c, &data.L, r.u)

	if data.L.UserID == 0 {
		// User isn't logged in
		// TODO: move this to a middleware somehow
		// TODO(ivan): status forbidden
		return errors.New("Must be logged in")
	}

	resp, err := r.s.ListArchivalEntries(context.TODO(), &schedulerproto.ListArchivalEntriesRequest{UserID: data.L.UserID})
	if err != nil {
		return err
	}

	data.ArchivalRequests = resp.Entries

	return c.JSON(http.StatusOK, data)
}
