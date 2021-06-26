package models

import (
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/jmoiron/sqlx"
)

// TODO: test suite for recommender implementations with precision and recall for sample dataset
type Recommender interface {
	GetRecommendations(userID int64) ([]int64, error)
}

// Dumb recommender system, computes expected rating value for user from a video's tags
// and orders by sum
// No more train otomads??? (please)
type BayesianTagSum struct {
	db *sqlx.DB
}

func NewBayesianTagSum(db *sqlx.DB) BayesianTagSum {
	return BayesianTagSum{
		db: db,
	}
}

func (b *BayesianTagSum) GetRecommendations(uid int64) ([]int64, error) {
	// Videos which have been viewed and not rated at implicitly 0
	sql := "WITH tag_ratings AS (select tag, sum(ratings.rating)/count(*) - 2.5 AS tag_score from videos INNER JOIN tags ON videos.id = tags.video_id LEFT JOIN ratings ON ratings.video_id = videos.id WHERE ratings.user_id = $1 AND videos.views > 0 GROUP BY tag), video_scores AS (SELECT videos.id, sum(tag_score) AS video_score from videos INNER JOIN tags ON tags.video_id = videos.id INNER JOIN tag_ratings ON tag_ratings.tag = tags.tag WHERE videos.views = 0 AND videos.transcoded = true GROUP BY videos.id) SELECT id from video_scores ORDER BY video_score DESC limit 10;"
	rows, err := v.db.Query(sql, a)
	if err != nil {
		return nil, err
	}

	var ret []uint64
	for rows.Next() {
		var vid uint64
		err = rows.Scan(&vid)
		if err != nil {
			return nil, err
		}
		ret = append(ret, vid)
	}

	return ret, nil
}
