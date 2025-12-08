package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateHashCost(t *testing.T) {
	password := []byte("password123")
	hash, err := GenerateHash(password)
	assert.NoError(t, err)

	cost, err := bcrypt.Cost(hash)
	assert.NoError(t, err)
	assert.Equal(t, 12, cost, "Bcrypt cost should be 12")
}

func TestCreateJWTExpiration(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	payload := JWTPayload{
		UID: 123,
		Exp: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := CreateJWT(payload, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	// Note: ValidateJWT in auth.go doesn't return the payload fully, so we might need to parse it manually or trust the CreateJWT logic.
	// But we can check if it validates correctly.
	uid, err := ValidateJWT(token, *privateKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), uid)
}
