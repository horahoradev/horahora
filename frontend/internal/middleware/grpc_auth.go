package middleware

import (
	"context"
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

func (j *JWTGRPCAuthenticator) GRPCAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: maybe add validation on the number of cookies

		if len(c.Cookies()) < 1 {
			return next(c)
		}

		jwt := c.Cookies()[0].String()

		uid, err := j.authenticate(jwt)
		if err != nil {
			log.Errorf("Error while authenticating: %s", err)
			return next(c)
		}

		c.Set(UserIDKey, uid)
		return next(c)
	}
}

func (j *JWTGRPCAuthenticator) authenticate(jwt string) (*int64, error) {
	jwtValidationRequest := &userproto.ValidateJWTRequest{
		Jwt: jwt,
	}

	// TODO: maybe add timeout
	validationResp, err := j.config.UserClient.ValidateJWT(context.Background(), jwtValidationRequest)
	if err != nil {
		return nil, err
	}

	if !validationResp.IsValid {
		return nil, errors.New("invalid jwt")
	}

	return &validationResp.Uid, nil
}
