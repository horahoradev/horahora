package routes

import (
	"context"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

func (v RouteHandler) handleSchedulerVideoApproavl(c echo.Context) error {
	id := c.Param("id")
	// idInt, err := strconv.ParseInt(id, 10, 64)
	// if err != nil {
	// 	return err
	// }

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank < 1 {
		// privileged user, can show unapproved videos
		// TODO(ivan): status forbidden
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	_, err = v.s.ApproveVideo(context.Background(), &schedulerproto.ApproveVideoReq{VideoID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (v RouteHandler) getUnapprovedVideos(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	resp, err := v.s.GetUnapprovedVideoList(context.TODO(), &schedulerproto.Empty{})
	if err != nil {
		return err
	}

	var videos []UnapprovedVideo
	for _, video := range resp.UnapprovedVideos {
		vid := UnapprovedVideo{
			VideoID: video.VideoID,
			URL:     video.Url,
		}
		videos = append(videos, vid)
	}

	return c.JSON(http.StatusOK, videos)
}
