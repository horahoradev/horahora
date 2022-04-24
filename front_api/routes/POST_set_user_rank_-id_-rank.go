package routes

import (
	"context"
	"net/http"
	"strconv"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

func (v RouteHandler) handleSetRank(c echo.Context) error {
	id := c.Param("userid")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	desiredRank := c.Param("rank")
	rankInt, err := strconv.ParseInt(desiredRank, 10, 64)
	if err != nil {
		return err
	}

	// Requesting profile
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	if profile.Rank != 2 {
		return c.String(http.StatusForbidden, "Insufficient user status")
	}

	req := userproto.SetRankRequest{
		UserID: idInt,
		Rank:   userproto.UserRank(rankInt),
	}

	_, err = v.u.SetUserRank(context.TODO(), &req)
	return err
}
