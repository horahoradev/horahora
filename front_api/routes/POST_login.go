package routes

import (
	"context"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

func (r RouteHandler) handleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO: grpc auth goes here
	loginReq := &userproto.LoginRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := r.u.Login(context.Background(), loginReq)
	if err != nil {
		return err
	}

	return setCookie(c, loginResp.Jwt)
}
