package schedule

import "C"
import (
	"context"
	"errors"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/horahoradev/horahora/scheduler/internal/models"
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
	return poller{Db: db, PollingDelay: time.Second * 15, Redsync: redsync}, nil
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
				log.Infof("Sending url %s with parent %s to be processed", item.URL, item.ParentURL)
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

	log.Info("Fetching categories")
	urls, err := p.getURLs()
	if err != nil {
		return nil, err
	}

	log.Info("Fetching videos to dl")

	// The rand offset is a bit of a hack to prevent video downloads from being attempted many times per video, resulting in many lock acquisition failures
	// TODO: improve

	var ret []*models.VideoDLRequest
	for _, url := range urls {
		sql := "SELECT v.id, v.video_id, v.url, downloads.id FROM downloads " +
			"INNER JOIN downloads_to_videos d ON downloads.id = d.download_id " +
			"INNER JOIN videos v ON d.video_id = v.id " +
			"WHERE downloads.url = $1 AND v.dlStatus = 0 " +
			"ORDER BY CHAR_LENGTH(v.video_ID) DESC, v.video_ID desc LIMIT 1 " +
			"OFFSET random() * LEAST(1000, (select count(*) from downloads INNER JOIN downloads_to_videos d ON downloads.id = d.download_id INNER JOIN videos v ON d.video_id = v.id  WHERE downloads.url = $1 AND v.dlStatus = 0)) "
		res, err := p.Db.Query(sql, url)
		if err != nil {
			return nil, err
		}

		for res.Next() {
			req := models.VideoDLRequest{
				ParentURL: url,
				Redsync: p.Redsync,
				Db:      p.Db,
			}

			err = res.Scan(&req.ID, &req.VideoID, &req.URL, &req.DownloaddID)
			if err != nil {
				return nil, err
			}

			log.Infof("Got %s", req.URL)
			ret = append(ret, &req)
		}

	}

	return ret, nil
}

func (p *poller) getURLs() ([]string, error) {
	// TODO: only select synced download categories
	sql := "select d.url,  count(user_id) * random() AS score FROM " +
		"user_download_subscriptions s " +
		"INNER JOIN downloads d ON d.id = s.download_id " +
		"WHERE d.id IN (select downloads.id from downloads INNER JOIN downloads_to_videos d ON downloads.id = d.download_id INNER JOIN videos v on d.video_id = v.id WHERE v.dlStatus = 0 GROUP BY downloads.id) " +
		"GROUP BY d.id ORDER BY score desc LIMIT 1"
	row := p.Db.QueryRow(sql)

	var url string
	var score float64

	err := row.Scan(&url, &score)
	if err != nil {
		return nil, err
	}

	return []string{
		url,
	}, nil
}
