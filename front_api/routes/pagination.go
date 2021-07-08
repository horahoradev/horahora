package routes

import "errors"

const (
	NumberOfPagesToDisplay = 5
	NumberOfVideosPerPage  = 50
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Page number starts at 1
func getPageRange(numberOfVideos, currentPage int) ([]int, error) {
	if currentPage < 1 {
		return nil, errors.New("current page must be at least 1")
	}

	startVideoNum := (currentPage - 1) * NumberOfVideosPerPage
	if startVideoNum >= numberOfVideos {
		return nil, errors.New("bad arguments: start video number is greater than or equal to the number of videos for the query")
	}

	startPage := max(currentPage-2, 1)
	var ret []int

	for currPage := startPage; (currPage-1)*NumberOfVideosPerPage < numberOfVideos && currPage < startPage+NumberOfPagesToDisplay; currPage++ {
		ret = append(ret, currPage)
	}

	return ret, nil
}
