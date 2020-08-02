package middleware

import (
	config "github.com/horahoradev/horahora/frontend/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGRPCAuth(t *testing.T) {
	cfg, err := config.New()
	assert.NoError(t, err)

	auth := NewGRPCAuth(cfg)

	jwt := "" // TODO: get JWT from userservice for test

	uid, err := auth.authenticate(jwt)
	assert.NoError(t, err)

	assert.NotNil(t, uid)
}
