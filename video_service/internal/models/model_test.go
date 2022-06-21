//go:build integration
// +build integration

package models

import (
	"context"
	proto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/horahoradev/horahora/video_service/internal/config"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var v *VideoModel

func init() {
	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	// Register user
	req := proto.RegisterRequest{
		Email:    "wow@wow.com",
		Username: "wow",
		Password: "securepassword",
	}

	for i := 0; i < 11; i++ {
		_, err = cfg.UserClient.Register(context.Background(), &req)
		if err != nil {
			log.Panic(err)
		}
	}

	v, err = NewVideoModel(cfg.SqlClient, cfg.UserClient)
	if err != nil {
		log.Panic(err)
	}

	_, err = v.SaveForeignVideo(context.Background(), "mytestvideo", "wow", "", "0",
		0, "", "testlocation", "newLoc", []string{"test"}, 10)
	if err != nil {
		log.Panic(err)
	}
}

func TestSaveForeignVideoAndObtainVideoList(t *testing.T) {

	videoList, err := v.GetVideoList(videoproto.SortDirection_desc, 0)
	assert.NoError(t, err)

	var videoListContainsVideo bool
	for _, val := range videoList {
		if val.VideoTitle == "mytestvideo" {
			videoListContainsVideo = true
		}
	}

	assert.Equal(t, true, videoListContainsVideo)
}

func TestRating(t *testing.T) {
	// Ensure that we can add ratings
	err := v.AddRatingToVideoID("0", "0", 10.0)
	assert.NoError(t, err)

	// Ensure that the video info reflects the rating that we added
	info, err := v.getBasicVideoInfo(10, "0")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, info.rating)
}

func TestViews(t *testing.T) {
	info, err := v.getBasicVideoInfo(10, "0")
	assert.NoError(t, err)
	originalViews := info.views

	err = v.IncrementViewsForVideo("0")
	assert.NoError(t, err)

	info, err = v.getBasicVideoInfo(10, "0")
	assert.NoError(t, err)
	assert.Equal(t, originalViews+1, info.views)
}
