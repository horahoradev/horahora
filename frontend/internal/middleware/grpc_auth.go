package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/horahoradev/horahora/frontend/internal/config"
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
const UserRank = "userrank"

func (j *JWTGRPCAuthenticator) GRPCAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: maybe add validation on the number of cookies

		if len(c.Cookies()) < 1 {
			return next(c)
		}

		jwt := c.Cookies()[0].Value

		jwtDecoded, err := base64.StdEncoding.DecodeString(jwt)
		if err != nil {
			log.Errorf("Failed to decode jwt. Err: %s", err)
			return next(c)
		}

		uid, err := j.authenticate(string(jwtDecoded))
		if err != nil {
			log.Errorf("Error while authenticating: %s", err)
			return next(c)
		}

		c.Set(UserIDKey, uid)

		// TODO: add other stuff in here like username, profile picture location, etc
		// If user is authenticated, get other metadata
		resp, err := j.config.UserClient.GetUserFromID(context.Background(), &userproto.GetUserFromIDRequest{UserID: uid})
		if err != nil {
			log.Errorf("Could not retrieve authenticated users metadata. Err: %s", err)
			return next(c)
		}

		c.Set(UserRank, int32(resp.Rank))
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
		return 0, err
	}

	if !validationResp.IsValid {
		return 0, errors.New("invalid jwt")
	}

	return validationResp.Uid, nil
}
