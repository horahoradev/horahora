package models

import (
	"context"

	"github.com/go-redsync/redsync"
	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
)

type contentType string

const (
	Tag      contentType = "tag"
	Channel  contentType = "channel"
	Playlist contentType = "playlist"
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

type ContentArchivalRequest struct {
	Website      proto.SupportedSite `db:"website"`
	ContentType  contentType         `db:"attribute_type"`  // "channel", "tag", or "playlist"
	ContentValue string              `db:"attribute_value"` // either the channel VideoID or the tag string
}

func (m *ArchiveRequestRepo) GetContentArchivalRequests(userID int64) ([]ContentArchivalRequest, error) {
	sql := "SELECT website, attribute_type, attribute_value FROM user_download_subscriptions s " +
		"INNER JOIN downloads d ON d.id = s.download_id WHERE user_id=$1"
	var archivalRequests []ContentArchivalRequest

	err := m.Db.Select(&archivalRequests, sql, userID)
	if err != nil {
		return nil, err
	}

	return archivalRequests, nil
}

func (m *ArchiveRequestRepo) New(contentType contentType, contentValue string, website proto.SupportedSite, userID int64) error {
	tx, err := m.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	var downloadID uint32
	row := tx.QueryRow("INSERT INTO downloads(date_created, website, attribute_type, attribute_value) "+
		"VALUES (Now(), $1, $2, $3) ON CONFLICT (website, attribute_type, attribute_value) "+
		"DO UPDATE set website = EXCLUDED.website RETURNING id", website, contentType, contentValue)
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
	res, err := m.Db.Query("SELECT id, website, attribute_type, attribute_value FROM downloads " +
		"WHERE last_synced IS NULL or last_synced + interval '1 day' * backoff_factor < Now()")
	if err != nil {
		return nil, err
	}

	var ret []CategoryDLRequest
	for res.Next() {
		c := CategoryDLRequest{Redsync: m.Redsync, Db: m.Db}

		err = res.Scan(&c.Id, &c.Website, &c.ContentType, &c.ContentValue)
		if err != nil {
			return nil, err
		}

		ret = append(ret, c)
	}

	return ret, nil
}
