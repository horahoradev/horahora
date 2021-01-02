package schedule

import (
	"context"
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
	return poller{Db: db, PollingDelay: time.Second * 60, Redsync: redsync}, nil
}

func (p *poller) PollDatabaseAndSendIntoQueue(ctx context.Context, videoQueue chan *models.VideoDLRequest) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Context done, returning from database poll loop")
			return nil

		default:
			itemsToSchedule, err := p.getVideos()
			if err != nil {
				if err != FailedToFetch {
					log.Errorf("failed to get items. Err: %s", err)
				}
				break // try again lol
			}

			for _, item := range itemsToSchedule {
				log.Infof("Sending %s %s %s to be processed", item.C.Website, item.C.ContentType, item.C.ContentValue)
				videoQueue <- item
			}
		}
		time.Sleep(p.PollingDelay)
	}

	return nil
}

var FailedToFetch = errors.New("failed to retrieve desired number of items")

func (p *poller) getVideos() ([]*models.VideoDLRequest, error) {
	// TODO: put this in a repo later

	categories, err := p.getCategories()
	if err != nil {
		return nil, err
	}

	var ret []*models.VideoDLRequest
	for _, category := range categories {
		sql := "SELECT v.id, v.video_id, v.url, downloads.id FROM downloads " +
			"INNER JOIN downloads_to_videos d ON downloads.id = d.download_id " +
			"INNER JOIN videos v ON d.video_id = v.id " +
			"WHERE downloads.website = $1 AND downloads.attribute_type = $2 AND downloads.attribute_value = $3 AND v.dlStatus = 0 LIMIT 10"
		res, err := p.Db.Query(sql, category.Website, category.ContentType, category.ContentValue)
		if err != nil {
			return nil, err
		}

		for res.Next() {
			req := models.VideoDLRequest{
				C:       category,
				Redsync: p.Redsync,
				Db:      p.Db,
			}

			err = res.Scan(&req.ID, &req.VideoID, &req.URL, &req.DownloaddID)
			if err != nil {
				return nil, err
			}
			ret = append(ret, &req)
			log.Info(ret)
		}

	}

	return ret, nil
}

func (p *poller) getCategories() ([]models.Category, error) {
	// TODO: implement me :)
	return []models.Category{
		{
			Website:      0,
			ContentType:  "tag",
			ContentValue: "音MAD",
		},
	}, nil
}

// dequeueFromDatabase pops the n most recent items from the database and timestamps them
// I'm using postgres as a message queue because it's easy
// requires isolation to be serial
//func (p *poller) dequeueFromDatabase(ctx context.Context, numItems int) ([]*models.CategoryDLRequest, error) {
//	tx, err := p.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
//	if err != nil {
//		return nil, err
//	}
//
//	rows, err := tx.Query("SELECT id, website, attribute_type, attribute_value FROM downloads" +
//		" WHERE lock < NOW() - INTERVAL '10 minutes' OR lock IS NULL ORDER BY last_polled DESC limit 1")
//	if err != nil {
//		return nil, err
//	}
//
//	var dlReqs []*models.CategoryDLRequest
//	// At this point, we've acquired the selected items
//	for rows.Next() {
//		i := models.CategoryDLRequest{}
//
//		err := rows.Scan(&i.VideoID, &i.Website, &i.ContentType, &i.ContentValue)
//		if err != nil {
//			return nil, err
//		}
//		i.Db = p.Db
//		i.Redsync = p.Redsync
//
//		dlReqs = append(dlReqs, &i)
//	}
//
//	if len(dlReqs) != numItems {
//		err := tx.Rollback()
//		if err != nil {
//			log.Error("Failed to rollback")
//		}
//		return nil, FailedToFetch
//	}
//
//	for _, req := range dlReqs {
//		results, err := tx.Exec("UPDATE downloads SET last_polled = NOW(), lock = NOW() WHERE id=$1", req.VideoID)
//		rowsAffected, err2 := results.RowsAffected()
//		if err2 != nil {
//			return nil, err2
//		}
//
//		if err != nil || rowsAffected < 0 {
//			log.Errorf("Failed to update with err %s. Rolling back...", err)
//			err2 := tx.Rollback()
//			if err2 != nil {
//				log.Errorf("Rollback failed! Err: %s", err2)
//			}
//			return nil, err
//		}
//	}
//
//	err = tx.Commit()
//	// TODO: do I need to rollback here?
//	if err != nil {
//		return nil, err
//	}
//
//	return dlReqs, nil
//}
