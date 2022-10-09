package models

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/horahoradev/horahora/user_service/errors"

	serror "errors"

	"google.golang.org/grpc/status"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
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
	db *sqlx.DB
	// TODO: do we really need a grpc client here? bad cohesion ;(
	grpcClient        proto.UserServiceClient
	ApprovalThreshold int
	r                 Recommender
}

func NewVideoModel(db *sqlx.DB, client proto.UserServiceClient, approvalThreshold int) (*VideoModel, error) {
	rec := NewBayesianTagSum(db)

	return &VideoModel{db: db,
		grpcClient: client,
		r:          &rec,
	}, nil
}

// check if user has been created
// if it hasn't, then create it
// list user as parent of this video
// FIXME this signature is too long lol
// If domesticAuthorID is 0, will interpret as foreign video from foreign user
func (v *VideoModel) SaveForeignVideo(ctx context.Context, title, description string, foreignAuthorUsername string, foreignAuthorID string,
	originalSite string, originalVideoLink, originalVideoID, newURI string, tags []string, domesticAuthorID int64) (int64, error) {
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	if domesticAuthorID == 0 && (foreignAuthorID == "" || foreignAuthorUsername == "") {
		return 0, serror.New("foreign author info cannot be blank")
	}

	if domesticAuthorID == 0 && (originalVideoLink == "" || originalVideoID == "") {
		return 0, serror.New("original video info cannot be blank")
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
			log.Infof("User %s does not exist for video %s, creating...", foreignAuthorUsername, originalVideoID)

			regReq := proto.RegisterRequest{
				Email:          "fake@user.com", // NO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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

func (v *VideoModel) ForeignVideoExists(foreignVideoID, website string) (bool, error) {
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

	err = v.r.RemoveRecommendedVideoForUser(ratingUID, videoID)
	if err != nil {
		log.Errorf("Failed to remove recommended video for user %d. Err: %s", ratingUID, err)
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
func (v *VideoModel) GetVideoList(direction videoproto.SortDirection, pageNum int64, fromUserID int64, searchVal string, showUnapproved bool,
	category videoproto.OrderCategory) ([]*videoproto.Video, error) {
	sql, err := v.generateVideoListSQL(direction, pageNum, fromUserID, searchVal, showUnapproved, category)
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
		video.AuthorID = basicInfo.authorID
		video.Views = uint64(views)

		// FIXME: nothing is quite as dumb as this
		// Need to remove the absolute path from mpd loc
		video.ThumbnailLoc = strings.Replace(mpdLoc, ".mpd", ".thumb", 1)

		// TODO: could alloc in advance
		results = append(results, &video)
	}

	return results, nil
}

// FIXME: optimization: move to redis once I figure out what types of queries are necessary
func (v *VideoModel) GetNumberOfSearchResultsForQuery(fromUserID int64, searchVal string, showUnapproved bool) (int64, error) {
	var sql string
	var args []interface{}
	switch {
	case fromUserID != 0:
		sql = "SELECT COUNT(*) FROM videos WHERE userID = $1 AND transcoded=true AND is_deleted=false"
		if !showUnapproved {
			// oh no no no no FIXME
			sql = sql + " AND is_approved = true"
		}

		args = []interface{}{fromUserID}
	case searchVal != "":
		// TODO: DRY
		dialect := goqu.Dialect("postgres")
		ds := dialect.
			Select(goqu.COUNT("videos.id")).
			From(
				goqu.T("videos"),
			)

		conditions := v.getConditions(extractSearchTerms(searchVal))

		ds = ds.
			Where(conditions...).
			Where(goqu.C("transcoded").Eq(true)).
			Where(goqu.C("is_deleted").Eq(false))
		var err error
		sql, _, err = ds.ToSQL()
		if err != nil {
			return 0, err
		}
		args = []interface{}{}
	default:
		sql = "SELECT COUNT(*) FROM videos WHERE transcoded=true AND is_deleted=false"
		if !showUnapproved {
			// oh no no no no FIXME
			sql = sql + " AND is_approved = true"
		}
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

func (v *VideoModel) generateVideoListSQL(direction videoproto.SortDirection, pageNum, fromUserID int64, searchVal string, showUnapproved bool, orderCategory videoproto.OrderCategory) (string, error) {
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

	case videoproto.OrderCategory_my_ratings:
		ds = ds.LeftJoin(
			goqu.T("ratings"),
			goqu.On(goqu.Ex{"videos.id": goqu.I("ratings.video_id")})).
			Where(goqu.I("user_id").Eq(fromUserID))
		switch direction {
		case videoproto.SortDirection_asc:
			ds = ds.Order(goqu.I("ratings.rating").Asc())
		case videoproto.SortDirection_desc:
			ds = ds.Order(goqu.I("ratings.rating").Desc())
		}
	}

	// Mutually exclusive for now, can change later if desired
	switch {
	case fromUserID != 0:
		ds = ds.
			Where(goqu.C("userid").Eq(fromUserID))
	case searchVal != "":
		conditions := v.getConditions(extractSearchTerms(searchVal))

		ds = ds.
			Where(conditions...)
	}

	if !showUnapproved {
		// only show approved
		ds = ds.Where(goqu.C("is_approved").Eq(true))
	}

	// Only show transcoded videos
	ds = ds.Where(goqu.C("transcoded").Eq(true))

	// Do not show deleted videos
	ds = ds.Where(goqu.C("is_deleted").Eq(false))

	// TODO: ensure that this is safe from sql injection
	// Maybe use prepared mode?
	// update: this should be fine.
	sql, _, err := ds.ToSQL()

	return sql, err
}

// This function is pretty horrifying, one of the worst things I've written
func (v *VideoModel) getConditions(include, exclude []string) []exp.Expression {
	queryCols := []string{"title", "tag", "username"}

	f := func(terms []string, include bool) []exp.Expression {
		var incQuery []exp.Expression

		for _, term := range terms {
			var currConds []exp.Expression
			for _, col := range queryCols {
				patt := term // Prefix matching only, but we could add a reversed b-tree index for fully fuzzy matching
				// Oh no
				if col == "tag" {
					t := goqu.Dialect("postgres").
						Select("videos.id").
						From(
							goqu.T("videos"),
						).Join(
						goqu.T("tags"),
						goqu.On(goqu.Ex{"videos.id": goqu.I("tags.video_id")})).
						Where(goqu.I("tag").Eq(patt))
					var exp exp.Expression
					if include {
						exp = goqu.I("videos.id").In(t)
					} else {
						exp = goqu.I("videos.id").NotIn(t)
					}
					currConds = append(currConds, exp)
					// Oh no no no
				} else if col == "username" {
					resp, err := v.grpcClient.GetUserIDsForUsername(context.Background(), &proto.GetUserIDsForUsernameRequest{
						Username: term + "%"})
					if err != nil || len(resp.UserIDs) == 0 {
						log.Errorf("could not retrieve user ids for username %s", term)
						continue
					}

					in := make([]interface{}, 0, len(resp.UserIDs))

					for _, id := range resp.UserIDs {
						in = append(in, id)
					}
					if include {
						exp := goqu.I("userid").In(in)
						currConds = append(currConds, exp)
					} else {
						exp := goqu.I("userID").NotIn(in)
						currConds = append(currConds, exp)
					}

				} else {
					if include {
						exp := goqu.I(col).Eq(patt)
						currConds = append(currConds, exp)
					} else {
						exp := goqu.I(col).Neq(patt)
						currConds = append(currConds, exp)
					}
				}
			}
			if include {
				incQuery = append(incQuery, goqu.Or(currConds...))
			} else {
				incQuery = append(incQuery, goqu.And(currConds...))

			}
		}
		return incQuery
	}

	// TODO: DRY
	incConds := f(include, true)
	excConds := f(exclude, false)
	allConds := append(incConds, excConds...)

	return []exp.Expression{goqu.And(allConds...)}
}

// TODO: might want to switch to some domain-specific information retrieval language in the future
// One of the log query languages, like Lucene, could work
func extractSearchTerms(search string) (includeTerms, excludeTerms []string) {
	var exclude, include []string
	spl := strings.Split(search, " ")
	for _, term := range spl {
		switch {
		case strings.HasPrefix(term, "-"):
			exclude = append(exclude, term[1:])
		default:
			include = append(include, term)
		}
	}
	return include, exclude
}

type basicVideoInfo struct {
	authorName string
	authorID   int64
	rating     float64
}

// Information that isn't super straightforward to query for
func (v *VideoModel) GetVideoInfo(videoID string) (*videoproto.VideoMetadata, error) {
	sql := "SELECT id, title, description, upload_date, userID, newLink, views FROM videos WHERE id=$1 AND is_deleted=false"
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
	videoInfo.authorID = authorID

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
	// Newest videos first
	sql := "SELECT id, newLink FROM videos WHERE transcoded = false AND too_big = false ORDER BY upload_date desc LIMIT 100"
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

func (v *VideoModel) MarkVideoAsTooBig(uv UnencodedVideo) error {
	sql := "UPDATE videos SET too_big = true WHERE id = $1"
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

func (v *VideoModel) GetVideoRecommendations(userID int64) (*videoproto.RecResp, error) {
	videoList, err := v.r.GetRecommendations(userID)
	if err != nil {
		log.Errorf("Could not get recommendations. Err: %s", err)
		return nil, err
	}

	return &videoproto.RecResp{Videos: videoList}, nil
}

func (v *VideoModel) DeleteVideo(videoID string) error {
	sql := "UPDATE videos SET is_deleted = true WHERE id = $1"
	_, err := v.db.Exec(sql, videoID)
	return err
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
		comment.AuthorProfileImageUrl = "/static/images/placeholder.png"

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
