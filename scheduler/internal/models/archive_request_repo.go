package models

import (
	"context"

	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
	"net/url"
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

func (m *ArchiveRequestRepo) GetContentArchivalRequests(userID int64) ([]string, error) {
	sql := "SELECT Url FROM user_download_subscriptions s " +
		"INNER JOIN downloads d ON d.id = s.download_id WHERE user_id=$1"
	var urls []string

	err := m.Db.Select(&urls, sql, userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
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
			return nil, err
		}

		ret = append(ret, c)
	}

	return ret, nil
}


func GetWebsiteFromURL(u string) (string, error) {
	urlParsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return urlParsed.Hostname(), nil
}