package schedule

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type VideoDlRequest struct {
	Website          string
	ContentType      string // "channel", "tag", or "playlist"
	ContentValue     string // either the channel ID or the tag string
	NumberToDownload int
	id               string
	db               *sqlx.DB
}

// RefreshLock refreshes the lock for this download request, preventing it from being acquired by another scheduler.
func (v *VideoDlRequest) RefreshLock() error {
	_, err := v.db.Exec("UPDATE downloads SET lock = Now() WHERE id = $1", v.id)
	return err
}

var NeverDownloaded error = errors.New("no video for category")

// Only relevant for tags
func (v *VideoDlRequest) GetLatestVideoForRequest() (*string, error) {
	curs, err := v.db.Query("SELECT video_ID from previous_downloads WHERE content_ID=$1 AND website=$2 ORDER BY upload_time desc LIMIT 1", v.ContentValue, v.Website)
	if err != nil {
		return nil, err
	}

	var videoIDList []string
	for curs.Next() {
		var i string
		curs.Scan(&i)
		videoIDList = append(videoIDList, i)
	}

	if len(videoIDList) == 0 {
		return nil, NeverDownloaded
	} else if len(videoIDList) != 1 {
		return nil, fmt.Errorf("videoIDList had the wrong length. Length: %d", len(videoIDList))
	}

	return &videoIDList[0], nil
}

func (v *VideoDlRequest) SetLatestVideo(videoID string, upload_time time.Time) error {
	_, err := v.db.Exec("INSERT INTO previous_downloads(video_ID, content_ID, upload_time) VALUES ($1, $2, $3)", videoID, v.ContentType, upload_time)
	return err
}
