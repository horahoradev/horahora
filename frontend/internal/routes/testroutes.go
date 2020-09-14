package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func SetupTestRoutes(e *echo.Echo) {
	e.GET("/", getTestHome)
	e.GET("/videos/:id", getTestVideo)
}

func getTestHome(c echo.Context) error {
	data := HomePageData{
		L: LoggedInUserData{},
		PaginationData: PaginationData{
			Pages:       []int{1, 2, 3, 4, 5},
			CurrentPage: 1,
			CurrentPath: "/",
		},
		Videos: []Video{
			{
				Title:        "[MAD] Barack Obama x スカイハイ",
				VideoID:      1,
				Views:        2,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       10.0,
			},
		},
	}

	return c.Render(http.StatusOK, "home", data)
}

func getTestVideo(c echo.Context) error {
	data := VideoDetail{
		L: LoggedInUserData{
			Username: "testuser",
			UserID:   1,
		},
		Title:           "testvideo",
		MPDLoc:          "",
		Views:           0,
		Rating:          10.0,
		AuthorID:        1, // TODO
		Username:        "testuser",
		UserDescription: "", // TODO: not implemented yet
		UserSubscribers: 0,  // TODO: not implemented yet
		ProfilePicture:  "/static/images/placeholder1.jpg",
		UploadDate:      time.Now().Format("2006-01-02"),
		VideoID:         1,
		Comments:        nil,
		Tags:            []string{"ytpmv", "test"},
	}

	return c.Render(http.StatusOK, "video", data)

}
