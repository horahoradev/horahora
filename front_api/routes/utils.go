package routes

import (
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
