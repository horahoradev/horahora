package routes

import (
	"context"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

// getArchiveRequests is a GET handler accepting no parameters, returning the list of archival entries
// Requires authentication
// route: GET /archiverequests
// Response is of this form:
// {"ArchivalRequests":[{"url":"https://www.youtube.com/watch?v=8DXqneHHzA8"}]}
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

	var requests []ArchivalRequest
	for _, archivalReq := range resp.Entries {
		req := ArchivalRequest{
			UserID:               archivalReq.UserID,
			Url:                  archivalReq.Url,
			ArchivedVideos:       archivalReq.ArchivedVideos,
			CurrentTotalVideos:   archivalReq.CurrentTotalVideos,
			LastSynced:           archivalReq.LastSynced,
			BackoffFactor:        archivalReq.BackoffFactor,
			DownloadID:           archivalReq.DownloadID,
			UndownloadableVideos: archivalReq.UndownloadableVideos,
		}
		requests = append(requests, req)
	}
	data.ArchivalRequests = requests

	return c.JSON(http.StatusOK, data)
}
