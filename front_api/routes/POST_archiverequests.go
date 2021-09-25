package routes

import (
	"context"
	"net/http"

	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) handleArchiveRequest(c echo.Context) error {
	urlVal := c.FormValue("url")

	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	req := schedulerproto.URLRequest{
		UserID: profile.UserID,
		Url: urlVal,
	}

	_, err = r.s.DlURL(context.TODO(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
