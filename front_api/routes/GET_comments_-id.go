package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (r RouteHandler) getComments(c echo.Context) error {
	videoID, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		return err
	}

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return err
	}

	// FIX LATER
	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64 for getComments")
	}

	resp, err := r.v.GetCommentsForVideo(context.Background(), &videoproto.CommentRequest{VideoID: videoIDInt, CurrUserID: UserIDInt})
	if err != nil {
		return err
	}

	commentList := make([]CommentData, 0)

	for _, comment := range resp.Comments {
		commentData := CommentData{
			ID:                 comment.CommentId,
			CreationDate:       comment.CreationDate,
			Content:            comment.Content,
			Username:           comment.AuthorUsername,
			ProfileImage:       comment.AuthorProfileImageUrl,
			VoteScore:          comment.VoteScore,
			CurrUserHasUpvoted: comment.CurrentUserHasUpvoted,
		}
		if comment.ParentId != 0 {
			commentData.ParentID = comment.ParentId
		}

		commentList = append(commentList, commentData)
	}

	return c.JSON(http.StatusOK, &commentList)
}
