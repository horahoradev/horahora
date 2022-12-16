package routes

import (
	"context"
	"net/http"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

const (
	minNameLength     = 5
	minPasswordLength = 5
)

// Route: POST /login
// Accepts form-encoded values: username, password
// response: 200 if ok, and sets a cookie
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
		return echo.NewHTTPError(http.StatusForbidden, "Login failed.")
	}

	return setCookie(c, loginResp.Jwt)
}
