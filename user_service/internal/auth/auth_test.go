package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordValidation(t *testing.T) {
	testPassword := "mystrongpassword"
	cost := 5
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(testPassword), cost)
	assert.NoError(t, err)

	isSame, err := compareHashedPassword([]byte(testPassword), hashedPass)
	assert.NoError(t, err)
	assert.Equal(t, isSame, true)
}

// Commands I used:
// openssl genrsa -out keypair.pem 2048
func TestKeypairParse(t *testing.T) {
	keypair := `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA+MHfP6bK1Tm5Qsy49WSD9TIsdKvstfnshIIlc6Or7jr4Lz/c
ZqivsGtsuenlyyMg0uEFKGC2VXojdyysPwsPN1V+OZNKrUexuUKCJ25xCob8xncJ
3ml2zLRyEKTAWTyV9KlaviapphoKdO+kUzVM+mM+BqXv2C/gS4nGhH4xN/8GRlGM
3x/w9d/kMbk0yVo8WvNCObAj1PBzTTc6pvXHijS8HwR4Kf6cUgEjLbAjDZqTAFQH
kL3wTTUspRbyCqNBY3pahqsMJsur3TnHBzdYrKhXOIdW5HF2H0Vi0URX03mD1uo3
775yQMy/R90pX2pYDXO50HVmlaBSRaRwVq4PWwIDAQABAoIBAQDOnDAx7PXxLfWK
3wEMiIT1hcvYx+deqgthb7ttlda6ku4KSI9mENuDu2Xd8MB0/WNI467l/5XR3xVP
6aHS5vunaPHjHkxbKf4aMnxqBdgL91xB9yMSQSR0F7fZzQ0Z0TF3uaXt46zqkhx3
TYd9uPqLyRv+Z5qTRAiWlQN2jl2Q9DRn3HdTaGPF5rsn2CjFPmoQGX7A+objSBSq
3JrU4cv6huqCZODqyCgVjdHfaRUFLGCPs5phoO29LY2DElo6BhWaJl2tvIXy2eqN
yqNwYjHHfMmzBv+OWyOY9IuussRajx4JvlIdh0M8Fm5Tniqeip/LpG0GOvK8l+qd
n3q3GLrJAoGBAP8SkhpL5A0HIkFwsL56eye9Ofquj71AoWeIxs3A/GYssjGoDQnU
rG3GUBFNuTSo5ulPSe43vJ3naTZZBffQXI8bW+2+dIAoExApu/jqH+ETfERhRQMS
4N8u88jxYTdi3WCOjhGpRjnXm9Wtwgq91wsvAAh+fv6b03FH1EZKgEDHAoGBAPmp
bEYK8QWg4AKK6e/EQhZW9su4Sdb38U4wQTyTvlDN7Bc+ZYvXeOlyBxdHDiEWC+XL
wHtrOlqQ+U9nCfC4lrJbGbVmrQ4mSwVhLGPQz2zntvxw8jQwZx+K4OFs4JFqHhWc
khpz21VGxWsJ3Jix6GoDKtz03IomcDpr9lxUZFDNAoGBAKJnMYz6qu28kAv4cyAk
Hcu0iHjasfw+bUXdaS7R5CIt7Rr+s6aBuXN/Y7VQtk4YCEWeTSUWacpj77JBxjH9
gSFAuyxJKiX63gBZgiw+7SNCY8mp4OXPHEwduexD+7DnCqqSuVP3YhYr+DV5l2V9
b7DYMP43hCYaEus6X6aNgtE7AoGBAIg6GgpeDgW0MocwpVVfEXB/I0sl06SoxdKU
IgSb2UzeD+Te9ynG+QLoZVYeP2duUC+jbfPqHn0sfd0FrDbdgdzwOKbyz5rY6jaV
P1N3rLcP+JjmSEKR5rMfZHWcoyy1apUASfiFHzj41OADEYuACAFQmSLXuT7omnRG
VLcslVBBAoGBAMJgmuUGIXvS4cBDsgjrKH7mCssjBandLN2NPkNN6uhSnANxm7n8
g36Dcu4E9TEcC7qNSJ8eVeutOENerGTVJ8fUAbXhIdvjcLz0iSMkpROCNz+Zr35G
rd258VnoYyNVswrjem4jHKTm4frORBF3sx6R1i/KiFSptp941g2hYjGe
-----END RSA PRIVATE KEY-----
`

	privateKey, err := ParsePrivateKey(keypair)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, privateKey)
}
