package routes

import (
	"context"
	"errors"
	"net/http"
	"strings"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /archiverequests
// Requires authentication
// Accepts form-encoded value URL, which is the url to be archived
// response: 200 if ok
func (r RouteHandler) handleArchiveRequest(c echo.Context) error {
	urlVal := strings.TrimSpace(c.FormValue("url"))

	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank <= 1 {
		return errors.New("insufficient permissions")
	}

	req := schedulerproto.URLRequest{
		UserID: profile.UserID,
		Url:    urlVal,
	}

	_, err = r.s.DlURL(context.TODO(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
