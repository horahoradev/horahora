package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"strings"

	"github.com/horahoradev/horahora/front_api/config"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type JWTGRPCAuthenticator struct {
	config *config.Config
}

func NewGRPCAuth(config *config.Config) *JWTGRPCAuthenticator {
	return &JWTGRPCAuthenticator{
		config: config,
	}
}

const UserIDKey = "userid"
const UserRankKey = "userrank"
const UserLoggedIn = "loggedin"

func (j *JWTGRPCAuthenticator) GRPCAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqPath := c.Request().URL.Path
		switch { // obv if they're trying to login or register we don't need to try to auth them
		case strings.HasPrefix(reqPath, "/api/login") || strings.HasPrefix(reqPath, "/api/register") || strings.HasPrefix(reqPath, "/api/logout"):
			return next(c)
		}

		c.Set(UserLoggedIn, false)

		if len(c.Cookies()) < 1 {
			return c.String(http.StatusForbidden, "")
		}

		jwt := c.Cookies()[0].Value

		jwtDecoded, err := base64.StdEncoding.DecodeString(jwt)
		if err != nil {
			log.Errorf("Failed to decode jwt. Err: %s", err)
			return c.Redirect(http.StatusMovedPermanently, "/authentication/login")
		}

		uid, err := j.authenticate(string(jwtDecoded))
		if err != nil {
			log.Errorf("Error while authenticating: %s", err)
			return c.Redirect(http.StatusMovedPermanently, "/authentication/login")
		}

		c.Set(UserIDKey, uid)

		// TODO: add other stuff in here like username, profile picture location, etc
		// If user is authenticated, get other metadata
		_, err = j.config.UserClient.GetUserFromID(context.Background(), &userproto.GetUserFromIDRequest{UserID: uid})
		if err != nil {
			log.Errorf("Could not retrieve authenticated users metadata. Err: %s", err)
			return c.String(http.StatusForbidden, "")
		}

		c.Set(UserLoggedIn, true)
		return next(c)
	}
}

func (j *JWTGRPCAuthenticator) authenticate(jwt string) (int64, error) {
	jwtValidationRequest := &userproto.ValidateJWTRequest{
		Jwt: jwt,
	}

	// TODO: maybe add timeout
	validationResp, err := j.config.UserClient.ValidateJWT(context.Background(), jwtValidationRequest)
	if err != nil {
		return 0, fmt.Errorf("could not validate jwt. Err: %s", err)
	}

	if !validationResp.IsValid {
		return 0, errors.New("invalid jwt")
	}

	return validationResp.Uid, nil
}
