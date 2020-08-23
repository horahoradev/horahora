package models

import (
	"fmt"
	proto "github.com/horahoradev/horahora/scheduler/protocol"
	"github.com/jmoiron/sqlx"
	"log"
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
	Website      Website     `db:"website"`
	ContentType  contentType `db:"attribute_type"`  // "channel", "tag", or "playlist"
	ContentValue string      `db:"attribute_value"` // either the channel ID or the tag string
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
	websiteVal, err := getWebsiteStringFromEnum(website)
	if err != nil {
		return err
	}

	_, err = m.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value, userID) "+
		"VALUES (Now(), $1, $2, $3, $4)", websiteVal, contentType, contentValue, userID)
	return err
}

type Website string

const (
	Niconico Website = "niconico"
	Bilibili Website = "bilibili"
	Youtube  Website = "youtube"
)

func (w Website) ToProtoSupportedSite() proto.SupportedSite {
	switch w {
	case Niconico:
		return proto.SupportedSite_niconico
	case Bilibili:
		return proto.SupportedSite_bilibili
	case Youtube:
		return proto.SupportedSite_youtube
	}

	// FIXME: will be removed once I get rid of all of these dumb types
	log.Fatalf("unsupported site %s", w)
	return proto.SupportedSite_niconico
}

// this is dumb but I don't know what's to be done about it...
func getWebsiteStringFromEnum(enumVal proto.SupportedSite) (Website, error) {
	switch enumVal {
	case proto.SupportedSite_niconico:
		return Niconico, nil
	case proto.SupportedSite_bilibili:
		return Bilibili, nil
	case proto.SupportedSite_youtube:
		return Youtube, nil
	default:
		return "", fmt.Errorf("could not find specified website for enum %d", enumVal)
	}
}
