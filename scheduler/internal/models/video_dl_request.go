package models

import (
	"errors"
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
	"time"
)

type VideoDlRequest struct {
	Website      Website
	ContentType  contentType // "channel", "tag", or "playlist"
	ContentValue string      // either the channel ID or the tag string
	Id           string
	Db           *sqlx.DB
	Redsync      *redsync.Redsync
}

func NewVideoDlRequest(website Website, contentType contentType, contentValue, id string, Db *sqlx.DB, redsync2 *redsync.Redsync) *VideoDlRequest {
	return &VideoDlRequest{
		Website:      website,
		ContentType:  contentType,
		ContentValue: contentValue,
		Id:           id,
		Db:           Db,
		Redsync:      redsync2,
	}
}

// RefreshLock refreshes the lock for this download request, preventing it from being acquired by another scheduler.
func (v *VideoDlRequest) RefreshLock() error {
	_, err := v.Db.Exec("UPDATE downloads SET lock = Now() WHERE id = $1", v.Id)
	return err
}

var NeverDownloaded error = errors.New("no video for category")

// Only relevant for tags
func (v *VideoDlRequest) GetLatestVideoForRequest() (*string, error) {
	curs, err := v.Db.Query("SELECT video_ID from previous_downloads WHERE content_ID=$1 AND website=$2 ORDER BY upload_time desc LIMIT 1", v.ContentValue, v.Website)
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
	_, err := v.Db.Exec("INSERT INTO previous_downloads(video_ID, content_ID, upload_time) VALUES ($1, $2, $3)", videoID, v.ContentType, upload_time)
	return err
}

func (v *VideoDlRequest) AcquireLockForVideo(videoID string) error {
	mut := v.Redsync.NewMutex(videoID, redsync.SetExpiry(time.Minute*30))
	return mut.Lock()
}
