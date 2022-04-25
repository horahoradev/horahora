package routes

import (
	"context"
	"net/http"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /comments/
// Requires authentication
// Accepts form-encoded values: video_id, content (content of comment), and parent (parent comment id if a reply)
// response: 200 if ok
func (r RouteHandler) handleComment(c echo.Context) error {
	videoID := c.FormValue("video_id")

	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	content := c.FormValue("content")

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return err
	}

	parent := c.FormValue("parent")
	parentIDInt, _ := strconv.ParseInt(parent, 10, 64)

	_, err = r.v.MakeComment(context.Background(), &videoproto.VideoComment{
		UserId:        profile.UserID,
		VideoId:       videoIDInt,
		Comment:       content,
		ParentComment: parentIDInt,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
