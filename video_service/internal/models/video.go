package models

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"git.horahora.org/otoman/user-service.git/errors"

	"google.golang.org/grpc/status"

	proto "git.horahora.org/otoman/user-service.git/protocol"

	_ "git.horahora.org/otoman/user-service.git/protocol"
	"github.com/jmoiron/sqlx"
)

type VideoModel struct {
	db         *sqlx.DB
	grpcClient proto.UserServiceClient
}

func NewVideoModel(db *sqlx.DB, client proto.UserServiceClient) (*VideoModel, error) {
	return &VideoModel{db: db,
		grpcClient: client}, nil
}

// check if user has been created
// if it hasn't, then create it
// list user as parent of this video
func (v *VideoModel) SaveForeignVideo(ctx context.Context, title, description string, authorUsername string, authorID string, originalSite proto.Site, originalVideoLink, newURI string, tags []string) (int64, error) {
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	req := proto.GetForeignUserRequest{
		OriginalWebsite: originalSite,
		ForeignUserID:   authorID,
	}

	var horahoraUID int64

	resp, err := v.grpcClient.GetUserForForeignUID(ctx, &req)
	grpcErr, ok := status.FromError(err)
	if !ok {
		return 0, fmt.Errorf("could not parse gRPC err")
	}
	switch {
	case grpcErr.Message() == errors.UserDoesNotExistMessage:
		// Create the user
		log.Info("User does not exist for video")

		regReq := proto.RegisterRequest{
			Email:          "",
			Username:       authorUsername,
			Password:       "",
			ForeignUser:    true,
			ForeignUserID:  authorID,
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

	sql := "INSERT INTO videos (title, description, userID, originalSite, " +
		"originalLink, newLink) " +
		"VALUES ($1, $2, $3, $4, $5, $6)" +
		"returning id"

	// By this point the user should exist
	// Username is unique, so will fail if user already exists
	var videoID int64
	res := tx.QueryRow(sql, title, description, horahoraUID, originalSite, originalVideoLink, newURI)

	err = res.Scan(&videoID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tagSQL := "INSERT INTO video_tags (video_id, tag) VALUES ($1, $2)"
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
