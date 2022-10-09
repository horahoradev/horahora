package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

// getArchiveRequests is a GET handler accepting no parameters, returning the list of archival entries
// Requires authentication
// route: GET /archiverequests
// Response is of this form:
// {"ArchivalRequests":[{"url":"https://www.youtube.com/watch?v=8DXqneHHzA8"}]}
func (r RouteHandler) getArchiveEvents(c echo.Context) error {
	downloadID, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if downloadID == "all" {
		data := ArchiveEventsData{}

		resp, err := r.s.ListArchivalEvents(context.TODO(), &schedulerproto.ListArchivalEventsRequest{DownloadID: 0, ShowAll: true})
		if err != nil {
			return err
		}

		data.ArchivalEvents = resp.Events
		return c.JSON(http.StatusOK, data)
	} else {
		downloadIDInt, err := strconv.ParseInt(downloadID, 10, 64)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		data := ArchiveEventsData{}

		resp, err := r.s.ListArchivalEvents(context.TODO(), &schedulerproto.ListArchivalEventsRequest{DownloadID: downloadIDInt})
		if err != nil {
			return err
		}

		data.ArchivalEvents = resp.Events

		return c.JSON(http.StatusOK, data)
	}

}
