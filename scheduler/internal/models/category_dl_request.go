package models

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type CategoryDLRequest struct {
	Url     string
	Website string
	Id      string
	Db      *sqlx.DB
}

// RefreshLock refreshes the lock for this download request, preventing it from being acquired by another scheduler.
func (v *CategoryDLRequest) RefreshLock() error {
	_, err := v.Db.Exec("UPDATE downloads SET lock = Now() WHERE id = $1", v.Id)
	return err
}

var NeverDownloaded = errors.New("no video for category")

const (
	MINIMUM_BACKOFF_TIME   = time.Hour * 24
	MAXIMUM_BACKOFF_FACTOR = 8 // 8 days
)

// IsBackingOff indicates whether the archive request is backing off from full content syncs
// Context: videos can be added to a category of content at any time; some categories are updated frequently, and some
// tend to be stagnant. We should vary the rate at which we fully sync content from a given category based on how
// frequently it's updated. Exponential backoff is used as the backoff strategy.
func (v *CategoryDLRequest) IsBackingOff() (bool, error) {
	var lastSynced time.Time
	var backoffFactor int
	sql := "SELECT last_synced, backoff_factor FROM downloads WHERE id = $1"

	rows := v.Db.QueryRow(sql, v.Id)
	err := rows.Scan(&lastSynced, &backoffFactor)
	if err != nil {
		return false, err
	}

	// I did what I had to do...
	return time.Now().Sub(lastSynced.Add(MINIMUM_BACKOFF_TIME*time.Duration(backoffFactor))) < 0, nil
}

func (v *CategoryDLRequest) ReportSyncHit() error {
	sql := "UPDATE downloads SET backoff_factor = 1, last_synced = Now() WHERE id = $1"
	_, err := v.Db.Exec(sql, v.Id)
	if err != nil {
		return err
	}

	return nil
}

func (v *CategoryDLRequest) ReportSyncMiss() error {
	// maybe there's an easier way to do this? It doesn't really matter though
	tx, err := v.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	var backoff_factor uint32
	row := tx.QueryRow("SELECT backoff_factor FROM downloads WHERE id = $1", v.Id)
	err = row.Scan(&backoff_factor)
	if err != nil {
		tx.Rollback()
		return err
	}

	sql := "UPDATE downloads SET backoff_factor = $1, last_synced = Now() WHERE id = $2"
	_, err = v.Db.Exec(sql, min(MAXIMUM_BACKOFF_FACTOR, backoff_factor*2), v.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Idempotent, ensures that videos are added and correct associations are created
// returns bool indicating whether something was added
func (v *CategoryDLRequest) AddVideo(videoID, url string) (bool, error) {
	tx, err := v.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return false, err
	}

	if url == "" {
		return false, errors.New("url cannot be blank")
	}

	if videoID == "" {
		return false, errors.New("video ID cannot be blank")
	}

	var id uint32
	sql := "INSERT INTO videos (video_ID, Url, website) VALUES ($1, $2, $3) " +
		"ON CONFLICT (video_ID, website) DO UPDATE set video_ID = EXCLUDED.video_ID RETURNING id"
	row := tx.QueryRow(sql, videoID, url, v.Website)
	err = row.Scan(&id)
	if err != nil {
		return false, err
	}

	sql = "INSERT INTO downloads_to_videos (download_id, video_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	res, err := tx.Exec(sql, v.Id, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return rowsAffected >= 1, nil
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
