package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) handleComment(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return err
	}

	data := c.Request().PostForm
	videoID, err := url.QueryUnescape(data.Get("video_id"))
	if err != nil {
		return err
	}

	userID, err := url.QueryUnescape(data.Get("user_id"))
	if err != nil {
		return err
	}

	content, err := url.QueryUnescape(data.Get("content"))
	if err != nil {
		return err
	}

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return err
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	parentIDInt, _ := getAsInt64(data, "parent")
	//if err != nil {
	//	// nothing
	//}

	_, err = r.v.MakeComment(context.Background(), &videoproto.VideoComment{
		UserId:        userIDInt,
		VideoId:       videoIDInt,
		Comment:       content,
		ParentComment: parentIDInt,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
