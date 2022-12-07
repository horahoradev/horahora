package routes

import (
	"context"
	"errors"
	"net/http"

	custommiddleware "github.com/horahoradev/horahora/front_api/middleware"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type LoggedInUserData struct {
	UserID            int64  `json:"userID"`
	Username          string `json:"username"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Email             string `json:"email"`
	Rank              int32  `json:"rank"`
	Banned            bool   `json:"banned"`
}

// Route: /currentuserprofile/
// Accepts no arguments
// Response is of form: {"userID":0,"username":"admin","profile_picture_url":"","rank":2}
func (v RouteHandler) getCurrentUserProfile(c echo.Context) error {
	profile, err := v.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, profile)
}

// TODO: this function is overcalled; we don't need to fetch profile data every time
// maybe just cache responses and call it a day
func (r *RouteHandler) getUserProfileInfo(c echo.Context) (*LoggedInUserData, error) {
	loggedIn := c.Get(custommiddleware.UserLoggedIn)

	loggedInBool, ok := loggedIn.(bool)
	if !ok || !loggedInBool {
		return nil, errors.New("User is not logged in")
	}

	l := LoggedInUserData{}

	id := c.Get(custommiddleware.UserIDKey)

	idInt, ok := id.(int64)
	if !ok {
		log.Error("Could not assert id to int64")
		return nil, errors.New("could not assert id to int64")
	}

	getUserReq := userproto.GetUserFromIDRequest{
		UserID: idInt,
	}

	userResp, err := r.u.GetUserFromID(context.TODO(), &getUserReq)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	l.Username = userResp.Username
	// l.ProfilePictureURL = userResp. // TODO
	l.Rank = int32(userResp.Rank)
	l.Email = userResp.Email
	l.UserID = idInt
	return &l, nil
}
