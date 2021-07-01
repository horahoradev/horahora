package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetupTestRoutes(e *echo.Echo) {
	e.GET("/", getTestHome)
	e.GET("/videos/:id", getTestVideo)
	e.GET("/login", getLogin)
	e.GET("/register", getRegister)
	e.GET("/comments/:id", getComments) // TODO: web client grpc
	e.GET("/upload", getUpload)

}

func getTestHome(c echo.Context) error {
	data := HomePageData{
		L: LoggedInUserData{},
		PaginationData: PaginationData{
			Pages:                []int{1, 2, 3, 4, 5}, // FIXME
			CurrentPage:          1,
			PathsAndQueryStrings: []string{"/", "/", "/", "/", "/"},
		},
		Videos: []Video{
			{
				Title:        "[MAD] Barack Obama x スカイハイ",
				VideoID:      1,
				Views:        2,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       5.0,
			},
			{
				Title:        "[MAD] Barack Obama x スカイハイ",
				VideoID:      2,
				Views:        2,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       5.0,
			},
			{
				Title:        "[MAD] Barack Obama x スカイハイ",
				VideoID:      3,
				Views:        2,
				AuthorID:     1,
				AuthorName:   "testuser",
				ThumbnailLoc: "/static/images/placeholder1.jpg",
				Rating:       5.0,
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
		Title:  "[MAD] Barack Obama x スカイハイ",
		MPDLoc: "",
		VideoDescription: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Views:           0,
		Rating:          10.0,
		AuthorID:        1, // TODO
		Username:        "testuser",
		UserDescription: "WOW this is my user description...", // TODO: not implemented yet
		UserSubscribers: 0,                                    // TODO: not implemented yet
		ProfilePicture:  "/static/images/placeholder1.jpg",
		UploadDate:      time.Now().Format("2006-01-02"),
		VideoID:         1,
		Comments:        nil,
		Tags:            []string{"ytpmv", "test"},
		RecommendedVideos: []Video{
			{
				Title:        "WOW",
				VideoID:      3,
				Views:        5,
				AuthorID:     10,
				AuthorName:   "TESTUSER",
				ThumbnailLoc: "loc",
				Rating:       5.0,
			},
		},
	}

	/*
		Title        string
		VideoID      int64
		Views        uint64
		AuthorID     int64
		AuthorName   string
		ThumbnailLoc string
		Rating       float64
	*/
	return c.Render(http.StatusOK, "video", data)

}

func getComments(c echo.Context) error {
	testComments := []CommentData{
		{
			1,
			"10-17-2020",
			"WOW",
			"testuser",
			"/static/images/placeholder1.jpg",
			5,
			true,
			0,
		},
		{
			2,
			"10-17-2020",
			"nice video.........",
			"testuser2",
			"/static/images/placeholder1.jpg",
			2,
			false,
			0,
		},
	}

	return c.JSON(http.StatusOK, &testComments)
}
