package routes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var Tests = []struct {
	NumVideos   int
	CurrentPage int
	Result      []int
}{
	{
		2363,
		5,
		[]int{3, 4, 5, 6, 7},
	},
	{
		100,
		1,
		[]int{1, 2},
	},

	{
		49,
		1,
		[]int{1},
	},
	{
		101,
		1,
		[]int{1, 2, 3},
	},
	{
		500,
		9,
		[]int{7, 8, 9, 10},
	},
	{
		500,
		10,
		[]int{8, 9, 10},
	},
}

func TestPaginationRange(t *testing.T) {
	for _, test := range Tests {
		pages, err := getPageRange(test.NumVideos, test.CurrentPage)
		assert.NoError(t, err)
		assert.Equal(t, test.Result, pages)
	}
}
