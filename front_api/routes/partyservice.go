package routes

import (
	"context"
	"net/http"
	"strconv"

	partyproto "github.com/horahoradev/horahora/partyservice/protocol"

	"github.com/labstack/echo/v4"
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

	url, err := url.Parse(videoURL))
	if err != nil {
		return nil, err
	}

	// happy path
	// FIXME
	videoID := path.Base(url.Path)
	log.Infof("Video ID: %v", videoID)

	resp, err := v.v.GetVideo(context.Background(), &videoservice.VideoRequest{
		VideoID: fmt.Sprintf("%v", videoID),
	})
	if err != nil {
		return nil, err
	}

	_, err = v.p.AddVideo(context.Background(), &partyproto.VideoRequest{
		PartyID:  partyID,
		VideoURL: videoURL,
	})
	if err != nil {
		return err
	}

	// Extremely dumb: just trigger a poll when we've added a new video
	// I will fix this later to include all relevant info in events
	// FIXME
	v.srv.BroadcastToRoom("", "bcast", "event:addvideo", "addvideo")

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
