package models

import (
	"fmt"
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

func NewArchiveRequest(db *sqlx.DB) *ArchiveRequestModel {
	return &ArchiveRequestModel{Db: db}
}

func (m *ArchiveRequestModel) New(contentType contentType, contentValue string, website proto.SupportedSite) error {
	websiteVal, err := getWebsiteStringFromEnum(website)
	if err != nil {
		return err
	}

	_, err = m.Db.Exec("INSERT INTO downloads(date_created, website, attribute_type, attribute_value) "+
		"VALUES (Now(), $1, $2, $3)", websiteVal, contentType, contentValue)
	return err
}

type Website string

const (
	Niconico Website = "niconico"
	Bilibili Website = "bilibili"
	Youtube  Website = "youtube"
)

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
