package routes

import (
	"context"
	"net/http"

	userproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

func (v RouteHandler) handleGetDownloadsInProgress(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	resp, err := v.s.GetDownloadsInProgress(context.TODO(), &userproto.DownloadsInProgressRequest{})
	if err != nil {
		return err
	}

	var videos []VideoInProgress
	for _, video := range resp.Videos {
		vid := VideoInProgress{
			Website:  video.Website,
			VideoID:  video.VideoID,
			DlStatus: video.DlStatus.String(),
		}
		videos = append(videos, vid)
	}

	return c.JSON(http.StatusOK, videos)
}
