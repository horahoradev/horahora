package model

import (
	proto "github.com/horahoradev/horahora/partyservice/protocol"
	"github.com/jmoiron/sqlx"
)

type PartyRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PartyRepo {
	return &PartyRepo{
		db: db,
	}
}

func (p *PartyRepo) CreateWatchParty(userID, channelID int) error {
	sql := "INSERT INTO parties (LeaderID, id) VALUES ($1, $2)"

	_, err := p.db.Exec(sql, userID, channelID)
	return err
}

func (p *PartyRepo) JoinWatchParty(userID, partyID int) error {
	sql := "INSERT INTO watchers (PartyID, UserID) VALUES ($1, $2)"

	_, err := p.db.Exec(sql, partyID, userID)
	return err
}

func (p *PartyRepo) DeleteFromWatchParty(userID, partyID int) error {
	sql := "DELETE FROM watchers where PartyID = $1 AND UserID = $2"
	_, err := p.db.Exec(sql, partyID, userID)
	return err
}

func (p *PartyRepo) NewVideo(videoLocation, Title string, VideoID, PartyID int) error {
	sql := "INSERT INTO video_queue (Title, PartyID, VideoID, Location, TS) VALUES ($1, $2, $3, $4, Now())"
	_, err := p.db.Exec(sql, Title, PartyID, VideoID, videoLocation)
	return err
}

func (p *PartyRepo) NextVideo(partyID int) error {
	sql := "DELETE FROM video_queue WHERE id in (select id from video_queue WHERE PartyID = $1 ORDER BY TS desc LIMIT 1 ) LIMIT 1"
	_, err := p.db.Exec(sql, partyID)
	return err
}

// Not needed for MVP
func (p *PartyRepo) UpdateLeader(partyID, userID int) error {
	return nil
	// TODO FIXME
	// sql := "UPDATE parties SET LeaderID = $1 WHERE "

	// _, err := p.db.Exec(sql, userID)
	// return err
}

func (p *PartyRepo) GetPartyState(partyID int) (*proto.PartyState, error) {
	var resp proto.PartyState
	sql := "SELECT UserID FROM watchers WHERE PartyID = $1"
	curs, err := p.db.Query(sql, partyID)
	if err != nil {
		return nil, err
	}

	for curs.Next() {
		var user proto.User
		err := curs.Scan(&user.UserID)
		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &user)
	}

	sql = "SELECT VideoID, Title, Location from video_queue WHERE PartyID = $1 ORDER BY TS asc"
	curs, err = p.db.Query(sql, partyID)
	if err != nil {
		return nil, err
	}

	for curs.Next() {
		var video proto.Video
		err := curs.Scan(&video.ID, &video.Title, &video.Location)
		if err != nil {
			return nil, err
		}

		resp.Videos = append(resp.Videos, &video)
	}

	return &resp, nil
}
