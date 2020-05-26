package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// This package is responsible for polling the database, and sending items into the channel

type VideoDlRequest struct {
	Website          string
	ContentType      string // "channel", "tag", or "playlist"
	ContentValue     string // either the channel ID or the tag string
	NumberToDownload int
	id               string
	db               *sqlx.DB
}

// RefreshLock refreshes the lock for this download request, preventing it from being acquired by another scheduler.
func (v *VideoDlRequest) RefreshLock() error {
	_, err := v.db.Exec("UPDATE downloads SET lock = Now() WHERE id = $1", v.id)
	return err
}

var NeverDownloaded error = errors.New("no video for category")

// Only relevant for tags
func (v *VideoDlRequest) GetLatestVideoForRequest() (*string, error) {
	curs, err := v.db.Query("SELECT video_ID from previous_downloads WHERE content_ID=$1 AND website=$2 ORDER BY upload_time desc LIMIT 1", v.ContentValue, v.Website)
	if err != nil {
		return nil, err
	}

	var videoIDList []string
	for curs.Next() {
		var i string
		curs.Scan(&i)
		videoIDList = append(videoIDList, i)
	}

	if len(videoIDList) == 0 {
		return nil, NeverDownloaded
	} else if len(videoIDList) != 1 {
		return nil, fmt.Errorf("videoIDList had the wrong length. Length: %d", len(videoIDList))
	}

	return &videoIDList[0], nil
}

func (v *VideoDlRequest) SetLatestVideo(videoID string, upload_time time.Time) error {
	_, err := v.db.Exec("INSERT INTO previous_downloads(video_ID, content_ID, upload_time) VALUES ($1, $2, $3)", videoID, v.ContentType, upload_time)
	return err
}

type poller struct {
	Db           *sqlx.DB
	PollingDelay time.Duration
}

func NewPoller(db *sqlx.DB) (poller, error) {
	return poller{Db: db, PollingDelay: time.Second * 5}, nil
}

func (p *poller) PollDatabaseAndSendIntoQueue(ctx context.Context, videoQueue chan *VideoDlRequest) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Context done, returning from database poll loop")
			return nil

		default:
			itemsToSchedule, err := p.dequeueFromDatabase(ctx, 1)
			if err != nil {
				if err != FailedToFetch {
					log.Errorf("failed to get items. Err: %s", err)
				}
				break // try again lol
			}

			for _, item := range itemsToSchedule {
				videoQueue <- item
			}
		}
		time.Sleep(p.PollingDelay)
	}

	return nil
}

var FailedToFetch error = errors.New("failed to retrieve desired number of items")

// dequeueFromDatabase pops the n most recent items from the database and timestamps them
// I'm using postgres as a message queue because it's easy
// requires isolation to be serial
func (p *poller) dequeueFromDatabase(ctx context.Context, numItems int) ([]*VideoDlRequest, error) {
	tx, err := p.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query("SELECT id, website, attribute_type, attribute_value FROM downloads " +
		"WHERE lock < NOW() - INTERVAL '30 minutes' OR lock IS NULL ORDER BY last_polled DESC limit 1")
	if err != nil {
		return nil, err
	}

	var dlReqs []*VideoDlRequest
	// At this point, we've acquired the selected items
	for rows.Next() {
		i := VideoDlRequest{}

		err := rows.Scan(&i.id, &i.Website, &i.ContentType, &i.ContentValue)
		if err != nil {
			return nil, err
		}
		i.db = p.Db
		i.NumberToDownload = 100 //  FIXME

		dlReqs = append(dlReqs, &i)
	}

	if len(dlReqs) != numItems {
		err := tx.Rollback()
		if err != nil {
			log.Error("Failed to rollback")
		}
		return nil, FailedToFetch
	}

	for _, req := range dlReqs {
		results, err := tx.Exec("UPDATE downloads SET last_polled = NOW(), lock = NOW() WHERE id=$1", req.id)
		rowsAffected, err2 := results.RowsAffected()
		if err2 != nil {
			return nil, err2
		}

		if err != nil || rowsAffected < 0 {
			log.Errorf("Failed to update with err %s. Rolling back...", err)
			err2 := tx.Rollback()
			if err2 != nil {
				log.Errorf("Rollback failed! Err: %s", err2)
			}
			return nil, err
		}
	}

	err = tx.Commit()
	// TODO: do I need to rollback here?
	if err != nil {
		return nil, err
	}

	return dlReqs, nil
}
