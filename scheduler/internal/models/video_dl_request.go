package models

import (
	"time"

	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
)

type Category struct {
	Website      int
	ContentType  string
	ContentValue string
}

type VideoDLRequest struct {
	Redsync     *redsync.Redsync
	Db          *sqlx.DB
	C           Category
	VideoID     string // Foreign ID
	ID          int    // Domestic ID
	DownloaddID int
	URL         string
}

func (v *VideoDLRequest) SetDownloaded() error {
	sql := "UPDATE videos SET download_id = $1, upload_time =  Now() WHERE video_ID = $2 AND id IN (select videos.id FROM videos " +
		"INNER JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id INNER JOIN downloads ON downloads_to_videos.download_id = downloads.id " +
		"WHERE videos.website = $3)"
	_, err := v.Db.Exec(sql, v.DownloaddID, v.ID, v.C.Website)
	if err != nil {
		return err
	}

	return v.SetDownloadSucceeded()
}

func (v *VideoDLRequest) SetDownloadSucceeded() error {
	sql := "UPDATE videos SET dlStatus = 1 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	return err
}

func (v *VideoDLRequest) SetDownloadFailed() error {
	sql := "UPDATE videos SET dlStatus = 2 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	return err
}

func (v *VideoDLRequest) AcquireLockForVideo() error {
	mut := v.Redsync.NewMutex(v.VideoID, redsync.SetExpiry(time.Minute*30))
	return mut.Lock()
}
