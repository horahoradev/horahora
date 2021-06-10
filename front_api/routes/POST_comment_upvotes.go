package routes

import (
	"context"
	"net/http"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) handleUpvote(c echo.Context) error {
	// DUMB!
	err := c.Request().ParseForm()
	if err != nil {
		return err
	}

	data := c.Request().PostForm

	commentID, err := getAsInt64(data, "comment_id")
	if err != nil {
		return err
	}

	userID, err := getAsInt64(data, "user_id")
	if err != nil {
		return err
	}

	hasUpvoted, err := getAsBool(data, "user_has_upvoted")
	if err != nil {
		return err
	}

	_, err = r.v.MakeCommentUpvote(context.Background(), &videoproto.CommentUpvote{
		CommentId: commentID,
		UserId:    userID,
		IsUpvote:  hasUpvoted,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
