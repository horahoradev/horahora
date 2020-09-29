package models

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/horahoradev/horahora/user_service/errors"

	"google.golang.org/grpc/status"

	proto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/jmoiron/sqlx"
)

const (
	maxRating         = 10.00
	NumResultsPerPage = 50
	cdnURL            = "images.horahora.org"
)

type VideoModel struct {
	db          *sqlx.DB
	grpcClient  proto.UserServiceClient
	redisClient *redis.Client
}

func NewVideoModel(db *sqlx.DB, client proto.UserServiceClient, redisClient *redis.Client) (*VideoModel, error) {
	return &VideoModel{db: db,
		grpcClient:  client,
		redisClient: redisClient}, nil
}

// check if user has been created
// if it hasn't, then create it
// list user as parent of this video
// FIXME this signature is too long lol
// If domesticAuthorID is 0, will interpret as foreign video from foreign user
func (v *VideoModel) SaveForeignVideo(ctx context.Context, title, description string, foreignAuthorUsername string, foreignAuthorID string,
	originalSite proto.Site, originalVideoLink, originalVideoID, newURI string, tags []string, domesticAuthorID int64) (int64, error) {
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	req := proto.GetForeignUserRequest{
		OriginalWebsite: originalSite,
		ForeignUserID:   foreignAuthorID,
	}

	horahoraUID := domesticAuthorID

	if horahoraUID == 0 {
		resp, err := v.grpcClient.GetUserForForeignUID(ctx, &req)
		grpcErr, ok := status.FromError(err)
		if !ok {
			return 0, fmt.Errorf("could not parse gRPC err")
		}
		switch {
		case grpcErr.Message() == errors.UserDoesNotExistMessage:
			// Create the user
			log.Info("User does not exist for video, creating...")

			regReq := proto.RegisterRequest{
				Email:          "",
				Username:       foreignAuthorUsername,
				Password:       "",
				ForeignUser:    true,
				ForeignUserID:  foreignAuthorID,
				ForeignWebsite: originalSite,
			}
			regResp, err := v.grpcClient.Register(ctx, &regReq)
			if err != nil {
				return 0, err
			}

			validateReq := proto.ValidateJWTRequest{
				Jwt: regResp.Jwt,
			}

			// The validation is superfluous, but we need the claims
			// FIXME: can probably optimize
			validateResp, err := v.grpcClient.ValidateJWT(ctx, &validateReq)
			if err != nil {
				return 0, err
			}

			if !validateResp.IsValid {
				return 0, fmt.Errorf("jwt invalid (this should never happen!)")
			}

			horahoraUID = validateResp.Uid

		case err != nil:
			return 0, err

		case err == nil:
			horahoraUID = resp.NewUID
		}
	}

	sql := "INSERT INTO videos (title, description, userID, originalSite, " +
		"originalLink, newLink, originalID, upload_date) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, Now())" +
		"returning id"

	// By this point the user should exist
	// Username is unique, so will fail if user already exists
	// FIXME: there might be some issues with error handling here. Should test to make sure scan returns ErrNoRows if insertion fail.
	// maybe switch to: https://github.com/jmoiron/sqlx/issues/154#issuecomment-148216948
	var videoID int64
	res := tx.QueryRow(sql, title, description, horahoraUID, originalSite, originalVideoLink, newURI, originalVideoID)

	err = res.Scan(&videoID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tagSQL := "INSERT INTO tags (video_id, tag) VALUES ($1, $2)"
	for _, tag := range tags {
		_, err = tx.Exec(tagSQL, videoID, tag)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		// What to do here? Rollback?
		return 0, err
	}

	return videoID, nil
}

func (v *VideoModel) ForeignVideoExists(foreignVideoID string, website videoproto.Website) (bool, error) {
	sql := "SELECT id FROM videos WHERE originalSite=$1 AND originalID=$2"
	var videoID int64
	res := v.db.QueryRow(sql, website, foreignVideoID)
	err := res.Scan(&videoID)
	switch {
	case err == sql2.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default: // err == nil
		return true, nil
	}
}

func (v *VideoModel) IncrementViewsForVideo(videoID string) error {
	// Sorted set with atomic incrementation
	// Every single command is atomic: https://www.slideshare.net/RedisLabs/atomicity-in-redis-thomas-hunter
	floatCmd := v.redisClient.HIncrBy(ViewHashNamespace, videoID, 1.00)
	if err := floatCmd.Err(); err != nil {
		return err
	}

	if err := v.ConstructSortedViewList(); err != nil {
		return err
	}

	return nil
}

// FIXME: clean up and document all of the crazy redis stuff I'm doing, maybe put into a single file
// Ensures that the video has been added to the views list in redis (starting at 0 views)
// To be used on video approval only
func (v *VideoModel) AssertViewsZero(videoID string) error {
	floatCmd := v.redisClient.HSet(ViewRankingNamespace, videoID, 0.00)
	return floatCmd.Err()
}

//func (v *VideoModel) AssertRatingsZero(videoID string) error {
//	floatCmd := v.redisClient.ZAdd(RatingNamespace, redis.Z{
//		Score:  0.00,
//		Member: videoID,
//	})
//	return floatCmd.Err()
//}

// FIXME: optimization. Switch to hash table for single video view fetches?
func (v *VideoModel) GetViewsForVideo(videoID string) (uint64, error) {
	// just fetch from sorted set
	floatCmd := v.redisClient.HGet(ViewHashNamespace, videoID)
	if floatCmd.Err() != nil {
		return 0, floatCmd.Err()
	}

	n, err := strconv.ParseInt(floatCmd.Val(), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint64(n), nil
}

// Ranking namespaces are for approved videos ONLY
const (
	ViewHashNamespace    = "videos:views"
	ViewRankingNamespace = "videos:viewranking"
	RatingNamespace      = "videos:ratings"
	ApprovalNamespace    = "videos:approvals"
)

func (v *VideoModel) GetVideosByViewsOrRatings(startNumber, endNumber int64, order videoproto.SortDirection, category videoproto.OrderCategory) ([]string, error) {
	var namespace string
	switch category {
	case videoproto.OrderCategory_views:
		namespace = ViewRankingNamespace
	case videoproto.OrderCategory_rating:
		namespace = RatingNamespace
	default:
		return nil, fmt.Errorf("unsupported category: %s", category.String())
	}

	switch order {
	case videoproto.SortDirection_asc:
		res := v.redisClient.ZRange(namespace, startNumber, endNumber)
		return res.Result()

	case videoproto.SortDirection_desc:
		res := v.redisClient.ZRevRange(namespace, startNumber, endNumber)
		return res.Result()
	}

	return nil, nil
}

func (v *VideoModel) AddRatingToVideoID(ratingUID, videoID string, ratingValue float64) error {
	// hash table for each video with key being user ID
	// really easy
	if ratingValue > 10.0 || ratingValue < 0.00 {
		return fmt.Errorf("invalid rating value: %f. Video ratings must be real numbers between 0 and 10.", ratingValue)
	}

	videoKey := fmt.Sprintf("ratings:%s", videoID)

	boolCmd := v.redisClient.HSet(videoKey, ratingUID, ratingValue)

	if err := boolCmd.Err(); err != nil {
		return err
	}

	// BIG FIXME
	if err := v.UpdateVideoRatingRankingsForAllViewedVideos(); err != nil {
		return err
	}

	return nil
}

// For now, this only supports either fromUserID or withTag. Can support both in future, need to switch to
// goqu and write better tests
func (v *VideoModel) GetVideoList(direction videoproto.SortDirection, pageNum int64, fromUserID int64, withTag string, showUnapproved bool) ([]*videoproto.Video, error) {
	sql, err := generateVideoListSQL(direction, pageNum, fromUserID, withTag, showUnapproved)
	if err != nil {
		return nil, err
	}

	var results []*videoproto.Video

	rows, err := v.db.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var video videoproto.Video
		var authorID int64
		var mpdLoc string
		err = rows.Scan(&video.VideoID, &video.VideoTitle, &authorID, &mpdLoc)
		if err != nil {
			return nil, err
		}

		basicInfo, err := v.getBasicVideoInfo(authorID, string(video.VideoID))
		if err != nil {
			return nil, err
		}

		video.Rating = basicInfo.rating
		video.AuthorName = basicInfo.authorName
		video.Views = basicInfo.views

		// FIXME: nothing is quite as dumb as this
		// Need to remove the absolute path from mpd loc
		video.ThumbnailLoc = strings.Replace(mpdLoc, ".mpd", ".png", 1)

		// TODO: could alloc in advance
		results = append(results, &video)
	}

	return results, nil
}

// FIXME: optimization: move to redis once I figure out what types of queries are necessary
func (v *VideoModel) GetNumberOfSearchResultsForQuery(fromUserID int64, withTag string) (int64, error) {
	var sql string
	var args []interface{}
	switch {
	case fromUserID != 0:
		sql = "SELECT COUNT(*) FROM videos WHERE userID = $1"
		args = []interface{}{fromUserID}
	case withTag != "":
		sql = "SELECT COUNT(DISTINCT video_id) FROM tags WHERE tag = $1"
		args = []interface{}{withTag}
	default:
		sql = "SELECT COUNT(*) FROM videos"
		args = []interface{}{}
	}

	rows, err := v.db.Query(sql, args...)
	if err != nil {
		return 0, err
	}

	var l int64
	for rows.Next() {
		err = rows.Scan(&l)
		if err != nil {
			return 0, err
		}
	}

	return l, nil
}

func generateVideoListSQL(direction videoproto.SortDirection, pageNum, fromUserID int64, withTag string, showUnapproved bool) (string, error) {
	minResultNum := (pageNum - 1) * NumResultsPerPage
	dialect := goqu.Dialect("postgres")

	ds := dialect.
		Select("videos.id", "title", "userid", "newlink").
		From(
			goqu.T("videos"),
		).
		Offset(uint(minResultNum)).
		Limit(NumResultsPerPage)

	switch direction {
	case videoproto.SortDirection_asc:
		ds = ds.Order(goqu.I("upload_date").Asc())
	case videoproto.SortDirection_desc:
		ds = ds.Order(goqu.I("upload_date").Desc())
	}

	// Mutually exclusive for now, can change later if desired
	switch {
	case fromUserID != 0:
		ds = ds.
			Where(goqu.C("userid").Eq(fromUserID))
	case withTag != "":
		ds = ds.Join(
			goqu.T("tags"),
			goqu.On(goqu.Ex{"videos.id": goqu.I("tags.video_id")})).
			Where(goqu.C("tag").Eq(withTag))
	}

	if !showUnapproved {
		// only show approved
		ds = ds.Where(goqu.C("is_approved").Eq(true))
	}

	// TODO: ensure that this is safe from sql injection
	// Maybe use prepared mode?
	sql, _, err := ds.ToSQL()

	return sql, err
}

type basicVideoInfo struct {
	views      uint64
	authorName string
	rating     float64
}

func (v *VideoModel) GetVideoInfo(videoID string) (*videoproto.VideoMetadata, error) {
	sql := "SELECT id, title, description, upload_date, userID, newLink FROM videos WHERE id=$1"
	var video videoproto.VideoMetadata
	var authorID int64

	row := v.db.QueryRow(sql, videoID)

	err := row.Scan(&video.VideoID, &video.VideoTitle, &video.Description, &video.UploadDate, &authorID, &video.VideoLoc)
	if err != nil {
		return nil, err
	}

	basicInfo, err := v.getBasicVideoInfo(authorID, string(video.VideoID))
	if err != nil {
		return nil, err
	}

	video.Rating = basicInfo.rating
	video.AuthorName = basicInfo.authorName
	video.Views = basicInfo.views
	video.AuthorID = authorID

	tags, err := v.getVideoTags(videoID)
	if err != nil {
		return nil, err
	}

	video.Tags = tags

	return &video, nil
}

type Tag struct {
	Tag string `db:"tag"`
}

func (v *VideoModel) getVideoTags(videoID string) ([]string, error) {
	sql := "SELECT tag from tags WHERE video_id = $1"
	var tags []Tag

	if err := v.db.Select(&tags, sql, videoID); err != nil {
		log.Errorf("Failed to retrieve video tags. Err: %s", err)
		return nil, err
	}

	var ret []string
	// FIXME
	for _, val := range tags {
		ret = append(ret, val.Tag)
	}

	return ret, nil
}

func (v *VideoModel) getBasicVideoInfo(authorID int64, videoID string) (*basicVideoInfo, error) {
	var videoInfo basicVideoInfo

	// Given user id, look up author name
	userReq := proto.GetUserFromIDRequest{
		UserID: authorID,
	}

	userResp, err := v.grpcClient.GetUserFromID(context.TODO(), &userReq)
	if err != nil {
		// maybe we should skip if we can't look them up?
		return nil, err
	}

	videoInfo.authorName = userResp.Username

	// Look up views from redis
	videoInfo.views, err = v.GetViewsForVideo(videoID)
	switch {
	case err != nil && err.Error() == "redis: nil":
		break
	case err != nil:
		return nil, err
	}

	// Look up ratings from redis
	videoInfo.rating, err = v.GetAverageRatingForVideoID(videoID)
	switch {
	case err != nil && err.Error() == "redis: nil":
		break
	case err != nil:
		return nil, err
	}

	return &videoInfo, nil
}

// FIXME: check for stupidity on a better day
// this is wonky but if an unapproved video isn't in the view ordered list, it wont be in the ranked list either
func (v *VideoModel) UpdateVideoRatingRankingsForAllViewedVideos() error {
	videos, err := v.GetVideosByViewsOrRatings(1, -1, videoproto.SortDirection_asc, videoproto.OrderCategory_views)
	if err != nil {
		return err
	}

	return v.UpdateVideoRatingRankings(videos)
}

// This is expensive and should only be done every so often
func (v *VideoModel) UpdateVideoRatingRankings(videoIDs []string) error {
	for _, videoID := range videoIDs {
		rating, err := v.GetAverageRatingForVideoID(videoID)
		if err != nil {
			return err
		}

		// Default to 0.00 if no ratings
		if math.IsNaN(rating) {
			rating = 0.00
		}

		// zadd sounds wrong but documentation says it's correct
		cmd := v.redisClient.ZAdd(RatingNamespace, redis.Z{Score: rating, Member: videoID})
		if err := cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (v *VideoModel) GetBasicSQLInfoForVideo(videoID string) (string, string, int64, error) {
	// newlink, title, userid
	var newlink, title string
	var authorID int64
	curs := v.db.QueryRow("select newlink, title, userID FROM videos WHERE id=$1", videoID)

	if err := curs.Scan(&newlink, &title, &authorID); err != nil {
		return "", "", 0, err
	}

	return newlink, title, authorID, nil
}

func (v *VideoModel) GetTopVideos(category videoproto.OrderCategory, order videoproto.SortDirection, startInd, endIndex int64) ([]*videoproto.Video, error) {
	var retList []*videoproto.Video

	videoIDs, err := v.GetVideosByViewsOrRatings(startInd, endIndex, order, category)
	if err != nil {
		return nil, err
	}

	for _, videoID := range videoIDs {
		idInt, err := strconv.ParseInt(videoID, 10, 64)
		if err != nil {
			return nil, err
		}

		newlink, title, authorID, err := v.GetBasicSQLInfoForVideo(videoID)
		if err != nil {
			log.Errorf("Could not obtain basic SQL info for video. Err: %s. Continuing...", err)
			continue
		}

		redisInfo, err := v.getBasicVideoInfo(authorID, videoID)
		if err != nil {
			log.Errorf("Could not obtain basic video info for video. Err: %s. Continuing...", err)
			continue
		}

		v := videoproto.Video{
			VideoTitle: title,
			Views:      redisInfo.views,
			Rating:     redisInfo.rating,
			// FIXME: I did it again...
			ThumbnailLoc: strings.Replace(newlink, ".mpd", ".png", 1),
			VideoID:      idInt,
			AuthorName:   redisInfo.authorName,
		}

		retList = append(retList, &v)
	}

	return retList, nil
}

func (v *VideoModel) GetAverageRatingForVideoID(videoID string) (float64, error) {
	// iterate through elements of hash table and compute the average
	// this is probably too expensive to do every time, so if it gets to be
	// an issue we can compute every ~30 mins and cache the result
	// alternatively could keep running total, probably doesn't matter
	// Idea: cache in sorted set with expiration time of 30 mins? can use to return sorted list to frontend
	ratingTotalNum := 0.00
	ratingTotalDenom := 0.00

	videoKey := fmt.Sprintf("ratings:%s", videoID)

	// according to docs, cursor value starts at 0, and server returns next value to pass in
	var cursorVal uint64 = 0

	scanIterator := v.redisClient.HScan(videoKey, cursorVal, "", 0).Iterator()

	// Every second element is a rating
	i := 0
	for scanIterator.Next() {
		log.Info(scanIterator.Val())
		if i%2 == 0 {
			i++
			continue
		}
		i++

		approved, err := v.IsVideoApproved(videoID)
		if err != nil {
			return 0.00, err
		}

		if !approved {
			// Don't add it to the list
			continue
		}

		rating, err := strconv.ParseFloat(scanIterator.Val(), 64)
		if err != nil {
			return 0.00, err
		}

		ratingTotalNum += rating
		ratingTotalDenom++
	}

	return ratingTotalNum / ratingTotalDenom, nil
}

// TODO: removing videos after approval. Maybe have an expiration?
// TODO: copy pasta
func (v *VideoModel) ConstructSortedViewList() error {
	// according to docs, cursor value starts at 0, and server returns next value to pass in
	var cursorVal uint64 = 0

	scanIterator := v.redisClient.HScan(ViewHashNamespace, cursorVal, "", 0).Iterator()

	// Every second element is a view
	i := 0
	var videoID string
	for scanIterator.Next() {
		log.Info(scanIterator.Val())
		if i%2 == 0 {
			videoID = scanIterator.Val()
			i++
			continue
		}
		i++

		// TODO: check if video approved
		approved, err := v.IsVideoApproved(videoID)
		if err != nil {
			return err
		}

		if !approved {
			// Don't add it to the list
			continue
		}

		view, err := strconv.ParseFloat(scanIterator.Val(), 64)
		if err != nil {
			return err
		}

		cmd := v.redisClient.ZAdd(ViewRankingNamespace, redis.Z{Score: view, Member: videoID})
		if err := cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (v *VideoModel) MarkVideoApproved(videoID string) error {
	cmd := v.redisClient.HSet(ApprovalNamespace, videoID, 1)
	if err := cmd.Err(); err != nil {
		return err
	}

	sql := "UPDATE videos SET is_approved = TRUE WHERE id = $1"
	_, err := v.db.Exec(sql, videoID)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoModel) IsVideoApproved(videoID string) (bool, error) {
	cmd := v.redisClient.HExists(ApprovalNamespace, videoID)
	return cmd.Val(), cmd.Err()
}

// TODO: refactor into separate models by concern
// e.g. approvalsmodel, viewmodel, etc

// Individual trusted user approves of the video
func (v *VideoModel) ApproveVideo(userID, videoID int) error {
	sql := "INSERT INTO approvals (user_id, video_id) VALUES ($1, $2)"
	_, err := v.db.Exec(sql, userID, videoID)
	if err != nil {
		return err
	}

	if err = v.MarkApprovals(); err != nil {
		return err
	}

	return nil
}

type Approval struct {
	videoID string `db:"video_id"`
}

func (v *VideoModel) MarkApprovals() error {
	var approvals []Approval
	sql := "SELECT video_id FROM approvals GROUP BY video_id HAVING count(*) > 3"

	if err := v.db.Select(&approvals, sql); err != nil {
		log.Errorf("Failed to retrieve video approvals. Err: %s", err)
		return err
	}

	for _, approval := range approvals {
		err := v.MarkVideoApproved(approval.videoID)
		if err != nil {
			// Log and move on
			log.Errorf("Could not mark video %s as approved. Err: %s", approval.videoID, err)
		}
	}

	return nil
}
