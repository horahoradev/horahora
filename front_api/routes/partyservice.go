package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	partyproto "github.com/horahoradev/horahora/partyservice/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

/*
   rpc NewWatchParty(NewPartyRequest) returns (NewPartyResponse) {}
   rpc BecomeLeader(PartyRequest) returns (LeaderResponse) {}
   rpc JoinParty(PartyRequest) returns (Empty) {}
   rpc HeartBeat(PartyRequest) returns (Empty) {}
   rpc GetPartyState(PartyRequest) returns (PartyState) {}

   rpc AddVideo (VideoRequest) returns (Empty) {}
   rpc NextVideo(PartyRequest) returns (Empty) {}
*/

/*
	POST /api/newwatchparty
	POST /api/joinwachparty
	POST /api/heartbeat
	GET /api/partystate:id
	POST /api/addvideo/:id
	POST /api/nextvideo/:id

*/
func (v RouteHandler) handleNewWatchParty(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	channelIDInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	_, err = v.p.NewWatchParty(context.Background(), &partyproto.NewPartyRequest{
		UserID:    profile.UserID,
		ChannelID: channelIDInt,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (v RouteHandler) handleJoinWatchParty(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	partyIDInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	_, err = v.p.JoinParty(context.Background(), &partyproto.PartyRequest{
		UserID:  profile.UserID,
		PartyID: partyIDInt,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (v RouteHandler) handleHeartbeat(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	partyID := c.FormValue("PartyID")
	partyIDInt, err := strconv.ParseInt(partyID, 10, 64)
	if err != nil {
		return err
	}

	_, err = v.p.HeartBeat(context.Background(), &partyproto.PartyRequest{
		UserID:  profile.UserID,
		PartyID: partyIDInt,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (v RouteHandler) handleGetPartyState(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}
	resp, err := v.p.GetPartyState(context.Background(), &partyproto.PartyRequest{
		UserID:  profile.UserID,
		PartyID: idInt,
	})
	/*
		title := c.FormValue("title")
		description := c.FormValue("description")
		tagsList := c.FormValue("tags")
	*/

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

type NewVideoMsg struct {
	Title    string
	ID       int
	Location string
}

func (v RouteHandler) handleAddVideo(c echo.Context) error {
	id := c.Param("id")
	partyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	_, err = v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	videoURL := c.FormValue("VideoURL")

	url, err := url.Parse(videoURL)
	if err != nil {
		return err
	}

	// happy path
	// FIXME
	videoID := path.Base(url.Path)
	log.Infof("Video ID: %v", videoID)

	resp, err := v.v.GetVideo(context.Background(), &videoproto.VideoRequest{
		VideoID: fmt.Sprintf("%v", videoID),
	})
	if err != nil {
		return err
	}

	_, err = v.p.AddVideo(context.Background(), &partyproto.VideoRequest{
		PartyID: partyID,
		Video: &partyproto.Video{
			Title:    resp.VideoTitle,
			ID:       resp.VideoID,
			Location: resp.VideoLoc,
		},
	})
	if err != nil {
		return err
	}
	msg := NewVideoMsg{
		Title:    resp.VideoTitle,
		ID:       int(resp.VideoID),
		Location: resp.VideoLoc,
	}

	v.srv.BroadcastToRoom("", "bcast", "event:addvideo", msg)

	return c.JSON(http.StatusOK, nil)
}

func (v RouteHandler) handleNextVideo(c echo.Context) error {
	id := c.Param("id")
	partyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	_, err = v.p.NextVideo(context.Background(), &partyproto.PartyRequest{
		PartyID: partyID,
		UserID:  profile.UserID,
	})
	if err != nil {
		return err
	}

	v.srv.BroadcastToRoom("", "bcast", "event:nextvideo", "nextvideo")

	return c.JSON(http.StatusOK, nil)
}
