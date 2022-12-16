package routes

import (
	"context"
	"net/http"
	"strings"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
)

// Route: POST /logout
// Accepts form-encoded values username, password, and email
// response: 200 if ok, and sets a cookie
func (r RouteHandler) handleRegister(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	registrationReq := userproto.RegisterRequest{
		Password: password,
		Username: username,
		Email:    email,
	}

	regisResp, err := r.u.Register(context.Background(), &registrationReq)
	switch {
	case err != nil && strings.Contains(err.Error(), "invalid"): // no, don't do this TODO FIXME
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case err != nil:
		return echo.NewHTTPError(http.StatusForbidden, "Registration failed.")
	}

	// NO!!! FIXME
	validateResp, err := r.u.ValidateJWT(context.TODO(), &userproto.ValidateJWTRequest{
		Jwt: regisResp.Jwt,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Authentication failed after registering; forward this to OtoMAN immediately (?!)")
	}

	_, err = r.u.AddAuditEvent(context.TODO(), &userproto.NewAuditEventRequest{
		Message: "New user has registered",
		User_ID: validateResp.Uid,
	})
	if err != nil {
		return err // If the audit event can't be created, fail the operation
	}

	// TODO: use registration JWT to auth

	return setCookie(c, regisResp.Jwt)
}
