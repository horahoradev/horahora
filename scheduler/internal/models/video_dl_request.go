package models

import (
	"github.com/go-redsync/redsync"
	"time"
)

func (v *VideoDlRequest) SetDownloaded(videoID string) error {
	sql := "UPDATE videos SET download_id = $1, upload_time =  Now() WHERE video_ID = $2 AND id IN (select videos.id FROM videos " +
		"INNER JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id INNER JOIN downloads ON downloads_to_videos.download_id = downloads.id " +
		"WHERE videos.website = $3)"
	_, err := v.Db.Exec(sql, v.Id, videoID, v.Website)
	return err
}

func (v *VideoDlRequest) AcquireLockForVideo(videoID string) error {
	mut := v.Redsync.NewMutex(videoID, redsync.SetExpiry(time.Minute*30))
	return mut.Lock()
}
