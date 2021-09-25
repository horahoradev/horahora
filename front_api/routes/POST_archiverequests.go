package routes

import (
	"context"
	"errors"
	"net/http"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	schedulerproto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (r RouteHandler) handleArchiveRequest(c echo.Context) error {
	urlVal := c.FormValue("url")

	userID := c.Get(custommiddleware.UserIDKey)
	UserIDInt, ok := userID.(int64)
	if !ok {
		log.Error("Could not assert userid to int64")
		return errors.New("could not assert userid to int64")
	}

	req := schedulerproto.URLRequest{
		UserID: UserIDInt,
		Url: urlVal,
	}

	_, err := r.s.DlURL(context.TODO(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
