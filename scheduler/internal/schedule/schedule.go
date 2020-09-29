package schedule

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redsync/redsync"
	"github.com/horahoradev/horahora/scheduler/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// This package is responsible for polling the database, and sending items into the channel

type poller struct {
	Db           *sqlx.DB
	PollingDelay time.Duration
	Redsync      *redsync.Redsync
}

func NewPoller(db *sqlx.DB, redsync *redsync.Redsync) (poller, error) {
	return poller{Db: db, PollingDelay: time.Second * 5, Redsync: redsync}, nil
}

func (p *poller) PollDatabaseAndSendIntoQueue(ctx context.Context, videoQueue chan *models.VideoDlRequest) error {
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
				log.Infof("Sending %s %s %s to be processed", item.Website, item.ContentType, item.ContentValue)
				videoQueue <- item
			}
		}
		time.Sleep(p.PollingDelay)
	}

	return nil
}

var FailedToFetch = errors.New("failed to retrieve desired number of items")

// dequeueFromDatabase pops the n most recent items from the database and timestamps them
// I'm using postgres as a message queue because it's easy
// requires isolation to be serial
func (p *poller) dequeueFromDatabase(ctx context.Context, numItems int) ([]*models.VideoDlRequest, error) {
	tx, err := p.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query("SELECT id, website, attribute_type, attribute_value FROM downloads " +
		"WHERE lock < NOW() - INTERVAL '30 minutes' OR lock IS NULL ORDER BY last_polled DESC limit 1")
	if err != nil {
		return nil, err
	}

	var dlReqs []*models.VideoDlRequest
	// At this point, we've acquired the selected items
	for rows.Next() {
		i := models.VideoDlRequest{}

		err := rows.Scan(&i.Id, &i.Website, &i.ContentType, &i.ContentValue)
		if err != nil {
			return nil, err
		}
		i.Db = p.Db
		i.Redsync = p.Redsync

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
		results, err := tx.Exec("UPDATE downloads SET last_polled = NOW(), lock = NOW() WHERE id=$1", req.Id)
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
