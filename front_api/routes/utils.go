package routes

import (
	"context"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"strconv"
)

func getPageNumber(c echo.Context) int64 {
	pageNumber := c.QueryParam("page")
	var pageNumberInt int64 = 1

	if pageNumber != "" {
		num, err := strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			log.Errorf("Invalid page number %s, defaulting to 1", pageNumber)
		}
		pageNumberInt = num
	}

	return pageNumberInt
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
