package routes

import (
	"context"
	"net/http"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /comment_upvotes
// Requires authentication
// Accepts form-encoded value comment_id, which is the url to be archived
// response: 200 if ok
func (r RouteHandler) handleUpvote(c echo.Context) error {

	commentID := c.FormValue("comment_id")

	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		return err
	}

	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	hasUpvoted := c.FormValue("user_has_upvoted")
	hasUpvotedBool, err := strconv.ParseBool(hasUpvoted)
	if err != nil {
		return err
	}

	_, err = r.v.MakeCommentUpvote(context.Background(), &videoproto.CommentUpvote{
		CommentId: commentIDInts,
		UserId:    profile.UserID,
		IsUpvote:  hasUpvotedBool,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
