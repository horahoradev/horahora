package routes

import (
	"context"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type LoggedInUserData struct {
	UserID            int64 `json:"userID"`
	Username          string `json:"username"`
	ProfilePictureURL string  `json:"profile_picture_url"`
	Rank              int32 `json:"rank"`
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
	l := LoggedInUserData{}

	//id := c.Get(custommiddleware.UserIDKey)
	//
	//idInt, ok := id.(int64)
	//if !ok {
	//	log.Error("Could not assert id to int64")
	//	return nil, errors.New("could not assert id to int64")
	//}

	getUserReq := userproto.GetUserFromIDRequest{
		UserID: 0,
	}

	userResp, err := r.u.GetUserFromID(context.TODO(), &getUserReq)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	l.Username = userResp.Username
	// l.ProfilePictureURL = userResp. // TODO
	l.UserID = 0
	l.Rank = int32(userResp.Rank)

	return &l, nil
}