package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) getComments(c echo.Context) error {
	videoID, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	videoIDInt, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	profileInfo, err := r.getUserProfileInfo(c)
	if err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	resp, err := r.v.GetCommentsForVideo(context.Background(), &videoproto.CommentRequest{VideoID: videoIDInt, CurrUserID: profileInfo.UserID})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
