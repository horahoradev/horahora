package models

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/horahoradev/horahora/user_service/errors"

	"google.golang.org/grpc/status"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/horahoradev/horahora/user_service/protocol"
	proto "github.com/horahoradev/horahora/user_service/protocol"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
)

const (
	maxRating         = 10.00
	NumResultsPerPage = 50
	cdnURL            = "images.horahora.org"
)

type VideoModel struct {
	db         *sqlx.DB
	grpcClient proto.UserServiceClient
	//redisClient *redis.Client
}

func NewVideoModel(db *sqlx.DB, client proto.UserServiceClient, redisClient *redis.Client) (*VideoModel, error) {
	return &VideoModel{db: db,
		grpcClient: client,
	}, nil
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

func (v *VideoModel) IncrementViewsForVideo(videoID int64) error {
	sql := "UPDATE videos SET views = views + 1 WHERE id = $1"
	_, err := v.db.Exec(sql, videoID)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoModel) AddRatingToVideoID(ratingUID, videoID int64, ratingValue float64) error {
	sql := "INSERT INTO ratings (user_id, video_id, rating) VALUES ($1, $2, $3)" +
		"ON CONFLICT (user_id, video_id) DO update SET rating = $4"
	_, err := v.db.Exec(sql, ratingUID, videoID, ratingValue, ratingValue)
	if err != nil {
		return err
	}

	rating, err := v.GetAverageRatingForVideoID(videoID)
	if err != nil {
		return err
	}

	_, err = v.db.Exec("UPDATE videos SET rating = $1 WHERE id = $2", rating, videoID)
	if err != nil {
		return err
	}

	return nil
}

// For now, this only supports either fromUserID or withTag. Can support both in future, need to switch to
// goqu and write better tests
func (v *VideoModel) GetVideoList(direction videoproto.SortDirection, pageNum int64, fromUserID int64, withTag string, showUnapproved bool,
	category videoproto.OrderCategory) ([]*videoproto.Video, error) {
	sql, err := generateVideoListSQL(direction, pageNum, fromUserID, withTag, showUnapproved, category)
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
		var authorID, views int64
		var mpdLoc string
		err = rows.Scan(&video.VideoID, &video.VideoTitle, &authorID, &mpdLoc, &views)
		if err != nil {
			return nil, err
		}

		basicInfo, err := v.getBasicVideoInfo(authorID, video.VideoID)
		if err != nil {
			return nil, err
		}

		video.Rating = basicInfo.rating
		video.AuthorName = basicInfo.authorName
		video.Views = uint64(views)

		// FIXME: nothing is quite as dumb as this
		// Need to remove the absolute path from mpd loc
		video.ThumbnailLoc = strings.Replace(mpdLoc, ".mpd", ".jpg", 1)

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

func generateVideoListSQL(direction videoproto.SortDirection, pageNum, fromUserID int64, withTag string, showUnapproved bool, orderCategory videoproto.OrderCategory) (string, error) {
	minResultNum := (pageNum - 1) * NumResultsPerPage
	dialect := goqu.Dialect("postgres")

	ds := dialect.
		Select("videos.id", "title", "userid", "newlink", "views").
		From(
			goqu.T("videos"),
		).
		Offset(uint(minResultNum)).
		Limit(NumResultsPerPage)

	switch orderCategory {
	case videoproto.OrderCategory_upload_date:
		switch direction {
		case videoproto.SortDirection_asc:
			ds = ds.Order(goqu.I("upload_date").Asc())
		case videoproto.SortDirection_desc:
			ds = ds.Order(goqu.I("upload_date").Desc())
		}

	case videoproto.OrderCategory_views:
		switch direction {
		case videoproto.SortDirection_asc:
			ds = ds.Order(goqu.I("views").Asc())
		case videoproto.SortDirection_desc:
			ds = ds.Order(goqu.I("views").Desc())
		}

	case videoproto.OrderCategory_rating:
		// could've done a WITH and inner join onto ratings table... but whatever, this is fine
		switch direction {
		case videoproto.SortDirection_asc:
			ds = ds.Order(goqu.I("rating").Asc())
		case videoproto.SortDirection_desc:
			ds = ds.Order(goqu.I("rating").Desc())
		}
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
	authorName string
	rating     float64
}

// Information that isn't super straightforward to query for
func (v *VideoModel) GetVideoInfo(videoID string) (*videoproto.VideoMetadata, error) {
	sql := "SELECT id, title, description, upload_date, userID, newLink, views FROM videos WHERE id=$1"
	var video videoproto.VideoMetadata
	var authorID, views int64

	row := v.db.QueryRow(sql, videoID)

	err := row.Scan(&video.VideoID, &video.VideoTitle, &video.Description, &video.UploadDate, &authorID, &video.VideoLoc, &views)
	if err != nil {
		return nil, err
	}

	basicInfo, err := v.getBasicVideoInfo(authorID, video.VideoID)
	if err != nil {
		return nil, err
	}

	video.Rating = basicInfo.rating
	video.AuthorName = basicInfo.authorName
	video.Views = uint64(views)
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

func (v *VideoModel) getUserInfo(authorID int64) (*userproto.UserResponse, error) {
	// Given user id, look up author name
	userReq := proto.GetUserFromIDRequest{
		UserID: authorID,
	}

	userResp, err := v.grpcClient.GetUserFromID(context.TODO(), &userReq)
	if err != nil {
		// maybe we should skip if we can't look them up?
		return nil, err
	}

	return userResp, nil
}

func (v *VideoModel) getBasicVideoInfo(authorID int64, videoID int64) (*basicVideoInfo, error) {
	var videoInfo basicVideoInfo

	var err error

	resp, err := v.getUserInfo(authorID)
	if err != nil {
		return nil, err
	}

	videoInfo.authorName = resp.Username

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

func (v *VideoModel) GetAverageRatingForVideoID(videoID int64) (float64, error) {
	var rating float64
	sql := "SELECT sum(rating)/count(*) FROM ratings WHERE video_id = $1"
	res := v.db.QueryRow(sql, videoID)
	if err := res.Scan(&rating); err != nil {
		// LMAO FIXME
		return 0.00, nil
	}

	return rating, nil
}

func (v *VideoModel) MarkVideoApproved(videoID string) error {
	sql := "UPDATE videos SET is_approved = TRUE WHERE id = $1"
	_, err := v.db.Exec(sql, videoID)
	if err != nil {
		return err
	}

	return nil
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
	VideoID string `db:"video_id"`
}

func (v *VideoModel) MarkApprovals() error {
	var approvals []Approval
	sql := "SELECT video_id FROM approvals GROUP BY video_id HAVING count(*) >= 1"

	if err := v.db.Select(&approvals, sql); err != nil {
		log.Errorf("Failed to retrieve video approvals. Err: %s", err)
		return err
	}

	for _, approval := range approvals {
		err := v.MarkVideoApproved(approval.VideoID)
		if err != nil {
			// Log and move on
			log.Errorf("Could not mark video %s as approved. Err: %s", approval.VideoID, err)
		}
	}

	return nil
}

// Comment stuff
func (v *VideoModel) MakeComment(userID, videoID, parentID int64, content string) error {
	switch parentID {
	case 0:
		sql := "INSERT INTO comments (user_id, video_id, comment, creation_date)" +
			" VALUES ($1, $2, $3, Now())"
		_, err := v.db.Exec(sql, userID, videoID, content)
		if err != nil {
			return err
		}

	default:
		sql := "INSERT INTO comments (user_id, video_id, parent_comment, comment, creation_date)" +
			" VALUES ($1, $2, $3, $4, Now())"
		_, err := v.db.Exec(sql, userID, videoID, parentID, content)
		if err != nil {
			return err
		}
	}

	return nil
}

type UnencodedVideo struct {
	ID      uint32 `db:"id"`
	NewLink string `db:"newlink"`
}

func (v UnencodedVideo) GetMPDUUID() string {
	spl := strings.Split(v.NewLink, "/")
	r := spl[len(spl)-1]
	return r[:len(r)-4]
}

func (v *VideoModel) GetUnencodedVideos() ([]UnencodedVideo, error) {
	sql := "SELECT id, newLink FROM videos WHERE transcoded = false"
	var videos []UnencodedVideo
	err := v.db.Select(&videos, sql)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (v *VideoModel) MarkVideoAsEncoded(uv UnencodedVideo) error {
	sql := "UPDATE videos SET transcoded = true WHERE id = $1"
	_, err := v.db.Exec(sql, uv.ID)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoModel) MakeUpvote(userID, commentID int64, isUpvote bool) error {
	// lol bruh moment
	// this is a hack... I'll probably keep it this way until we have actual downvotes though
	var voteScore = 0
	if isUpvote {
		voteScore = 1
	}

	sql := "INSERT INTO comment_upvotes (user_id, comment_id, vote_score) VALUES ($1, $2, $3)" +
		"ON CONFLICT (user_id, comment_id) DO update SET vote_score = $4"
	_, err := v.db.Exec(sql, userID, commentID, voteScore, voteScore)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoModel) GetComments(videoID, currUserID int64) ([]*videoproto.Comment, error) {
	var comments []*videoproto.Comment
	sql := "SELECT id, sum(COALESCE(vote_score, 0)) as upvote_score, comments.user_id," +
		" creation_date, comment, COALESCE(parent_comment, 0) " +
		"FROM comments LEFT JOIN comment_upvotes ON id = comment_id GROUP BY id,comment_upvotes.comment_id HAVING video_id = $1"
	rows, err := v.db.Query(sql, videoID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment videoproto.Comment

		err = rows.Scan(&comment.CommentId, &comment.VoteScore, &comment.AuthorId,
			&comment.CreationDate, &comment.Content, &comment.ParentId)
		if err != nil {
			log.Errorf("Failed to scan. Err: %s", err)
			continue
		}
		resp, err := v.getUserInfo(comment.AuthorId)
		if err != nil {
			log.Errorf("Failed to retrieve username for comment with author id %d. Err: %s",
				comment.AuthorId, err)
			continue
		}

		comment.AuthorUsername = resp.Username
		comment.AuthorProfileImageUrl = "/static/images/placeholder1.jpg"

		var score uint64
		sqlTwo := "SELECT vote_score FROM comment_upvotes WHERE user_id = $1 AND comment_id = $2"

		res := v.db.QueryRow(sqlTwo, currUserID, comment.CommentId)
		err = res.Scan(&score)
		if err == nil && score > 0 {
			comment.CurrentUserHasUpvoted = true
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}
