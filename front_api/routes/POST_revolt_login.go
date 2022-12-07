package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type revoltPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RouteHandler) handleRevoltLogin(c echo.Context) error {
	// Are we logged in?
	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	log.Errorf("LOGIN FOR %v", profile.Email)
	payload := revoltPayload{
		Email:    profile.Email,
		Password: "null01010",
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://api:8000/auth/session/login", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// if resp.StatusCode >= 400 {
	// 	return fmt.Errorf("bad revolt response status for registration: %v", resp.StatusCode)
	// }

	// or something...
	// TODO: audit
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Errorf("SMEGMA: %s", string(bodyResp))

	return c.String(http.StatusOK, string(bodyResp))
}
