package routes

import (
	"context"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Route: POST /login
// Accepts form-encoded values: username, password
// response: 200 if ok, and sets a cookie
func (r RouteHandler) handlePasswordReset(c echo.Context) error {
	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	oldPass := c.FormValue("old_password")
	newPass := c.FormValue("new_password")

	log.Error(oldPass, newPass)

	// TODO: grpc auth goes here
	resetReq := userproto.ResetPasswordRequest{
		UserID:      profile.UserID,
		OldPassword: oldPass,
		NewPassword: newPass,
	}

	_, err = r.u.ResetPassword(context.TODO(), &resetReq)
	return err

}
