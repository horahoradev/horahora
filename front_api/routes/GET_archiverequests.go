package routes

import (
	"context"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

// getArchiveRequests is a GET handler accepting no parameters, returning the list of archival entries
// Response is of this form:
//
func (r RouteHandler) getArchiveRequests(c echo.Context) error {
	data := ArchiveRequestsPageData{}

	profileInfo, err := r.getUserProfileInfo(c)
	if err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	resp, err := r.s.ListArchivalEntries(context.TODO(), &schedulerproto.ListArchivalEntriesRequest{UserID: profileInfo.UserID})
	if err != nil {
		return err
	}

	data.ArchivalRequests = resp.Entries

	return c.JSON(http.StatusOK, data)
}
