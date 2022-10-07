package models

import (
	"context"
	"database/sql"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

// ArchiveRequest is the model for the creation of new archive requests
// It's very similar to video_dl_request, but video_dl_request has different utility.
type ArchiveRequestRepo struct {
	Db *sqlx.DB
}

// FIXME: this API feels a little dumb

func NewArchiveRequest(db *sqlx.DB) *ArchiveRequestRepo {
	return &ArchiveRequestRepo{Db: db}
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

func (m *ArchiveRequestRepo) GetArchivalEvents(downloadID int64, showAll bool) ([]Event, error) {
	var events []Event

	var rows *sql.Rows
	var err error

	switch showAll {
	case true:
		sql := "Select video_url, parent_url, event_message, event_time FROM archival_events ORDER BY event_time DESC LIMIT 100"

		rows, err = m.Db.Query(sql)
		if err != nil {
			return nil, err
		}

	case false:
		sql := "Select video_url, parent_url, event_message, event_time FROM archival_events WHERE download_id = $1 ORDER BY event_time DESC LIMIT 100"

		rows, err = m.Db.Query(sql, downloadID)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var event Event

		err = rows.Scan(&event.VideoURL, &event.ParentURL, &event.Message, &event.EventTimestamp)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (m *ArchiveRequestRepo) GetContentArchivalRequests(userID int64) ([]Archival, error) {
	// This query is pretty dumb. Event and archive queries should be separate. FIXME
	// TODO: this query should be joined on archival subscriptions, not the download user id
	// This is an MVP fix
	// nvm i misread it is lol
	sql := "SELECT Url, coalesce(last_synced, Now()), backoff_factor, downloads.id FROM " +
		"downloads INNER JOIN user_download_subscriptions s ON downloads.id = s.download_id WHERE s.user_id=$1"

	rows, err := m.Db.Query(sql, userID)
	if err != nil {
		return nil, err
	}

	var archives []Archival
	urlMap := make(map[string]bool)

	for rows.Next() {
		var archive Archival

		err = rows.Scan(&archive.Url, &archive.LastSynced, &archive.BackoffFactor, &archive.DownloadID)
		if err != nil {
			return nil, err
		}

		// ok so this is really dumb but I'm going to let it go for now.
		// This prevents duplicates
		// FIXME
		_, ok := urlMap[archive.Url]
		if !ok {
			progressSql := "WITH numerator AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id WHERE downloads_to_videos.download_id = $1 AND dlstatus>0 AND dlstatus <3), " +
				"denominator AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id  WHERE downloads_to_videos.download_id = $1), " +
				"undownloadable AS (select count(*) from videos LEFT JOIN downloads_to_videos ON videos.id = downloads_to_videos.video_id  WHERE downloads_to_videos.download_id = $1 AND videos.dlstatus=2) " +
				"SELECT (select * from numerator) as numerator, (select * from denominator) as denominator, (select * from undownloadable) AS undownloadable"

			row := m.Db.QueryRow(progressSql, archive.DownloadID)
			if err != nil {
				return nil, err
			}

			err = row.Scan(&archive.Numerator, &archive.Denominator, &archive.Undownloadable)
			if err != nil {
				return nil, err
			}

			urlMap[archive.Url] = true
			archives = append(archives, archive)
		}

	}

	return archives, nil
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
	res, err := m.Db.Query("SELECT downloads.id, Url FROM downloads INNER JOIN user_download_subscriptions s ON downloads.id = s.download_id " +
		"WHERE last_synced IS NULL or last_synced + interval '1 day' * backoff_factor < Now() GROUP BY downloads.id")
	if err != nil {
		return nil, err
	}

	var ret []CategoryDLRequest
	for res.Next() {
		c := CategoryDLRequest{Db: m.Db}

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
	sql := "UPDATE videos SET dlstatus = 0 WHERE videos.id IN (SELECT videos.id FROM videos INNER JOIN downloads_to_videos ON downloads_to_videos.video_id = videos.id INNER JOIN user_download_subscriptions ON downloads_to_videos.download_id = user_download_subscriptions.download_id WHERE user_download_subscriptions.download_id = $1 AND user_download_subscriptions.user_id = $2 AND dlstatus = 2)"
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
	sql := "select video_id, website, dlStatus FROM videos WHERE dlstatus >= 3 ORDER BY dlStatus ASC"

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
