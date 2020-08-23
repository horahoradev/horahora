package models

import (
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
type ArchiveRequestModel struct {
	Db *sqlx.DB
}

// FIXME: this API feels a little dumb

func NewArchiveRequest(db *sqlx.DB) *ArchiveRequestModel {
	return &ArchiveRequestModel{Db: db}
}

type ContentArchivalRequest struct {
	Website      proto.SupportedSite `db:"website"`
	ContentType  contentType         `db:"attribute_type"`  // "channel", "tag", or "playlist"
	ContentValue string              `db:"attribute_value"` // either the channel ID or the tag string
}

func (m *ArchiveRequestModel) GetContentArchivalRequests(userID int64) ([]ContentArchivalRequest, error) {
	sql := "SELECT website, attribute_type, attribute_value FROM downloads WHERE userID=$1"
	var archivalRequests []ContentArchivalRequest

	err := m.Db.Select(&archivalRequests, sql, userID)
	if err != nil {
		return nil, err
	}

	return archivalRequests, nil
}

func (m *ArchiveRequestModel) New(contentType contentType, contentValue string, website proto.SupportedSite, userID int64) error {

	_, err := m.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value, userID) "+
		"VALUES (Now(), $1, $2, $3, $4)", website, contentType, contentValue, userID)
	return err
}
