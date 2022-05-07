package models

import (
	"errors"
	"strings"
	"sync"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
)

// TODO: test suite for recommender implementations with precision and recall for sample dataset
type Recommender interface {
	GetRecommendations(userID int64) ([]*videoproto.VideoRec, error)
	RemoveRecommendedVideoForUser(userID, videoID int64) error
}

// Dumb recommender system, computes expected rating value for user from a video's tags
// and orders by sum
// No more train otomads??? (please)
type BayesianTagSum struct {
	db            *sqlx.DB
	storedResults map[int64][]*videoproto.VideoRec
	mut           sync.Mutex
}

func NewBayesianTagSum(db *sqlx.DB) BayesianTagSum {
	return BayesianTagSum{
		db:            db,
		storedResults: make(map[int64][]*videoproto.VideoRec),
	}
}

func (b *BayesianTagSum) GetRecommendations(uid int64) ([]*videoproto.VideoRec, error) {
	// Is there a cached item?
	// TODO: rw lock
	b.mut.Lock()
	items, ok := b.storedResults[uid]
	b.mut.Unlock()
	switch {
	case ok && len(items) != 0: // if we have a nonempty hit, return that
		return items, nil
	default:
		// otherwise, compute new results and store them
		results, err := b.getRecommendations(uid)
		if err != nil {
			return nil, err
		}
		b.mut.Lock()
		b.storedResults[uid] = results
		b.mut.Unlock()
		return results, nil
	}
}

func (b *BayesianTagSum) RemoveRecommendedVideoForUser(userID, videoID int64) error {
	b.mut.Lock()
	videos, ok := b.storedResults[userID]
	b.mut.Unlock()
	if !ok {
		return errors.New("No videos for given user")
	}

	for i, video := range videos {
		if video.VideoID == videoID {
			// I don't like the look of this. FIXME
			b.mut.Lock()
			b.storedResults[userID] = append(videos[:i], videos[i+1:]...)
			b.mut.Unlock()
			return nil
		}
	}

	return errors.New("Desired video could not be removed, was not found")
}

func (b *BayesianTagSum) getRecommendations(uid int64) ([]*videoproto.VideoRec, error) {
	// Videos which have been viewed and not rated are implicitly rated 0
	// left join from video scores returns some random videos by default
	sql := "WITH tag_ratings AS (select tag, coalesce(avg(ratings.rating), 0.00) AS tag_score from ratings INNER JOIN tags ON ratings.video_id = tags.video_id WHERE ratings.user_id = $1 GROUP BY tag), " +
		"video_scores AS (SELECT tags.video_id, coalesce(avg(tag_score), 0.00) AS video_score from  tags INNER JOIN tag_ratings ON tag_ratings.tag = tags.tag WHERE tag_score >= 3.5 GROUP BY tags.video_id ORDER BY video_score DESC, tags.video_id LIMIT 50) " +
		"SELECT videos.id, title, newLink from video_scores INNER JOIN videos ON video_scores.video_id = videos.id WHERE videos.is_deleted IS false AND videos.transcoded IS true AND videos.id NOT IN (SELECT video_id FROM ratings WHERE ratings.user_id = $1) limit 10"
	rows, err := b.db.Query(sql, uid)
	if err != nil {
		return nil, err
	}

	var ret []*videoproto.VideoRec
	for rows.Next() {
		vid := videoproto.VideoRec{}
		err = rows.Scan(&vid.VideoID, &vid.VideoTitle, &vid.ThumbnailLoc)
		if err != nil {
			return nil, err
		}

		// I should stop doing this...
		vid.ThumbnailLoc = strings.Replace(vid.ThumbnailLoc, ".mpd", ".thumb", 1)

		ret = append(ret, &vid)
	}

	return ret, nil
}
