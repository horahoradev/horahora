package routes

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/gommon/log"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
)

// route: GET /videos/:id where id is the video id
// The video id is located within the path. No other parameters are accepted.
// Response is of the form:{"Title":"コダック","MPDLoc":"http://localhost:9000/otomads/207f773c-1e23-11ec-a6c3-0242ac1c0004.mpd","Views":2,"Rating":0,"VideoID":5,"AuthorID":5,"Username":"たっぴ","UserDescription":"","VideoDescription":"YouTube　\u003ca href=\"https://youtu.be/kP_lYd9D2to\" target=\"_blank\" rel=\"noopener nofollow\"\u003ehttps://youtu.be/kP_lYd9D2to\u003c/a\u003e","UserSubscribers":0,"ProfilePicture":"/static/images/placeholder1.jpg","UploadDate":"2021-09-25T17:07:56.400857Z","Comments":null,"Tags":null}
// authorID, userDescription, and userSubscribers all have no meaning as of yet.
func (v *RouteHandler) getVideo(c echo.Context) error {
	id := c.Param("id")

	// Dumb
	videoID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	// Increment views first
	viewReq := videoproto.VideoViewing{VideoID: videoID}
	_, err = v.v.ViewVideo(context.Background(), &viewReq)
	if err != nil {
		return err
	}

	videoReq := videoproto.VideoRequest{
		VideoID: id,
	}

	videoInfo, err := v.v.GetVideo(context.Background(), &videoReq)
	if err != nil {
		return err
	}

	rating := videoInfo.Rating

	// lol
	if math.IsNaN(rating) {
		rating = 0.00
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		log.Errorf("failed to load user profile. Err: %s", err)
	}

	recResp, err := v.v.GetVideoRecommendations(context.Background(), &videoproto.RecReq{
		UserId: profile.UserID,
	})
	if err != nil {
		log.Errorf("Could not retrieve recommendations. Err: %s", err)
		// Continue anyway
	}

	recVideos := []Video{}
	if recResp != nil {
		for _, rec := range recResp.Videos {
			// FIXME: fill other fields after modifying protocol
			vid := Video{
				Title:        rec.VideoTitle,
				VideoID:      rec.VideoID,
				ThumbnailLoc: rec.ThumbnailLoc,
			}

			recVideos = append(recVideos, vid)
		}
	}

	data := VideoDetail{
		Title:             videoInfo.VideoTitle,
		MPDLoc:            videoInfo.VideoLoc, // FIXME: fix this in videoservice LOL this is embarrassing
		Views:             videoInfo.Views,
		Rating:            rating,
		AuthorID:          videoInfo.AuthorID, // TODO
		Username:          videoInfo.AuthorName,
		UserDescription:   "", // TODO: not implemented yet
		VideoDescription:  videoInfo.Description,
		UserSubscribers:   0, // TODO: not implemented yet
		ProfilePicture:    "/static/images/placeholder1.jpg",
		UploadDate:        videoInfo.UploadDate,
		VideoID:           videoInfo.VideoID,
		Tags:              videoInfo.Tags,
		RecommendedVideos: recVideos,
		L:                 profile,
	}

	return c.JSON(http.StatusOK, data)
}
