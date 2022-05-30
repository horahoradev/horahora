package models

import (
	"context"
	"math"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
)

// ArchiveRequest is the model for the creation of new archive requests
// It's very similar to video_dl_request, but video_dl_request has different utility.
type ArchiveRequestRepo struct {
	Db      *sqlx.DB
	Redsync *redsync.Redsync
}

// FIXME: this API feels a little dumb

func NewArchiveRequest(db *sqlx.DB, rs *redsync.Redsync) *ArchiveRequestRepo {
	return &ArchiveRequestRepo{Db: db,
		Redsync: rs}
}

type Archival struct {
	Url            string
	DownloadID     uint64
	Denominator    uint64
	Numerator      uint64
	Undownloadable uint64
	LastSynced     string
	BackoffFactor  uint32
}

type Event struct {
	ParentURL      string
	VideoURL       string
	Message        string
	EventTimestamp string
}

func (m *ArchiveRequestRepo) GetContentArchivalRequests(userID int64) ([]Archival, []Event, error) {
	// This query is pretty dumb. Event and archive queries should be separate. FIXME
	// TODO: this query should be joined on archival subscriptions, not the download user id
	// This is an MVP fix
	// nvm i misread it is lol
	sql := "SELECT Url, coalesce(last_synced, Now()), backoff_factor, downloads.id, coalesce(archival_events.video_url, ''), coalesce(archival_events.parent_url, ''), coalesce(event_message, ''), coalesce(event_time, Now()) FROM " +
		"downloads INNER JOIN user_download_subscriptions s ON downloads.id = s.download_id LEFT JOIN archival_events ON downloads.id = archival_events.download_id WHERE s.user_id=$1 ORDER BY event_time DESC"

	rows, err := m.Db.Query(sql, userID)
	if err != nil {
		return nil, nil, err
	}

	var archives []Archival
	var events []Event
	urlMap := make(map[string]bool)

	for rows.Next() {
		var archive Archival
		var event Event

		err = rows.Scan(&archive.Url, &archive.LastSynced, &archive.BackoffFactor, &archive.DownloadID, &event.VideoURL,
			&event.ParentURL, &event.Message, &event.EventTimestamp)
		if err != nil {
			return nil, nil, err
		}

		if event.ParentURL != "" {
			events = append(events, event)
		}

		// ok so this is really dumb but I'm going to let it go for now.
		// This prevents duplicates
		// FIXME
		_, ok := urlMap[archive.Url]
		if !ok {
			progressSql := "WITH numerator AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id WHERE downloads_to_videos.download_id = $1 AND dlstatus>0 AND dlStatus <3), " +
				"denominator AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id  WHERE downloads_to_videos.download_id = $1), " +
				"undownloadable AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id  WHERE downloads_to_videos.download_id = $1 AND videos.dlStatus=2) " +
				"SELECT (select * from numerator) as numerator, (select * from denominator) as denominator, (select * from undownloadable) AS undownloadable"

			row := m.Db.QueryRow(progressSql, archive.DownloadID)
			if err != nil {
				return nil, nil, err
			}

			err = row.Scan(&archive.Numerator, &archive.Denominator, &archive.Undownloadable)
			if err != nil {
				return nil, nil, err
			}

			urlMap[archive.Url] = true
			archives = append(archives, archive)
		}

	}

	// This slice is similarly dumb
	// also a minor FIXME
	if events != nil {
		events = events[:uint64(math.Min(200, float64(len(events))))]
	}

	return archives, events, nil
}

func (m *ArchiveRequestRepo) New(url string, userID int64) error {
	tx, err := m.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	var downloadID uint32
	row := tx.QueryRow("INSERT INTO downloads(date_created, Url) "+
		"VALUES (Now(), $1) ON CONFLICT (Url) "+
		"DO UPDATE set Url = EXCLUDED.Url RETURNING id", url)
	err = row.Scan(&downloadID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO user_download_subscriptions (user_id, download_id) VALUES ($1, $2)", userID, downloadID)
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

func (m *ArchiveRequestRepo) GetUnsyncedCategoryDLRequests() ([]CategoryDLRequest, error) {
	// Fetch all unsynced downloads, irrespective of priority
	res, err := m.Db.Query("SELECT id, Url FROM downloads " +
		"WHERE last_synced IS NULL or last_synced + interval '1 day' * backoff_factor < Now()")
	if err != nil {
		return nil, err
	}

	var ret []CategoryDLRequest
	for res.Next() {
		c := CategoryDLRequest{Redsync: m.Redsync, Db: m.Db}

		err = res.Scan(&c.Id, &c.Url)
		if err != nil {
			return nil, err
		}

		c.Website, err = GetWebsiteFromURL(c.Url)
		if err != nil {
			log.Errorf("Failed to parse %s", c.Url)
			continue
		}

		ret = append(ret, c)
	}

	return ret, nil
}

func (m *ArchiveRequestRepo) DeleteArchivalRequest(userID, downloadID uint64) error {
	sql := "DELETE FROM user_download_subscriptions WHERE user_id = $1 AND download_id = $2"

	_, err := m.Db.Exec(sql, userID, downloadID)
	return err
}

func (m *ArchiveRequestRepo) RetryArchivalRequest(userID, downloadID uint64) error {
	sql := "UPDATE videos SET dlStatus = 0 WHERE videos.id IN (SELECT videos.id FROM videos INNER JOIN downloads_to_videos ON downloads_to_videos.video_id = videos.id INNER JOIN user_download_subscriptions ON downloads_to_videos.download_id = user_download_subscriptions.download_id WHERE user_download_subscriptions.download_id = $1 AND user_download_subscriptions.user_id = $2 AND dlStatus = 2)"
	_, err := m.Db.Exec(sql, downloadID, userID)
	return err
}

type Video struct {
	ID       string `db:"video_id"`
	Website  string `db:"website"`
	DlStatus int    `db:"dlstatus"` // postgres lowercased it, lol
}

func (m *ArchiveRequestRepo) GetDownloadsInProgress() ([]Video, error) {
	var videos []Video
	sql := "select video_id, website, dlStatus FROM videos WHERE dlStatus >= 3"
	err := m.Db.Select(&videos, sql)
	return videos, err
}

func (m *ArchiveRequestRepo) WipeDownloadsInProgress() error {
	sql := "UPDATE videos SET dlStatus = 0 WHERE dlStatus >= 3"
	_, err := m.Db.Exec(sql)
	return err
}

func GetWebsiteFromURL(u string) (string, error) {
	urlParsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return urlParsed.Hostname(), nil
}
