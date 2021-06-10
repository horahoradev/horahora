package routes

import (
	"context"
	"errors"
	"net/http"

	custommiddleware "github.com/SEAPUNK/horahora/front_api/middleware"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (r RouteHandler) handleArchiveRequest(c echo.Context) error {
	website := c.FormValue("website")
	contentType := c.FormValue("contentType")
	contentValue := c.FormValue("contentValue")

	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	websiteEnumVal, ok := schedulerproto.SupportedSite_value[website]
	if !ok {
		return errors.New("site not found")
	}

	supportedWebsite := schedulerproto.SupportedSite(websiteEnumVal)

	// FIXME: this is dumb. Fix this to use schedulerproto consts after switching to string instead of enum
	switch contentType {
	case "tag":
		req := schedulerproto.TagRequest{
			UserID:   UserIDInt,
			Website:  supportedWebsite, // FIXME: placeholder, see above
			TagValue: contentValue,
		}

		_, err := r.s.DlTag(context.TODO(), &req)
		if err != nil {
			return err
		}
	case "channel":
		req := schedulerproto.ChannelRequest{
			Website:   supportedWebsite,
			ChannelID: contentValue,
		}

		_, err := r.s.DlChannel(context.TODO(), &req)
		if err != nil {
			return err
		}
	case "playlist":
		req := schedulerproto.PlaylistRequest{
			Website:    supportedWebsite,
			PlaylistID: contentValue,
		}
		_, err := r.s.DlPlaylist(context.TODO(), &req)
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid content type")
	}

	return c.JSON(http.StatusOK, nil)
}
