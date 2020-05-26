package model

import (
	"database/sql"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/horahoradev/horahora/user_service/errors"

	dbsql "database/sql"

	proto "github.com/horahoradev/horahora/user_service/protocol"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserModel struct {
	Conn *sqlx.DB
}

func NewUserModel(db *sqlx.DB) (*UserModel, error) {
	u := &UserModel{Conn: db}

	return u, nil
}

// Password in this context is in plaintext
func (m *UserModel) NewUser(username, email string, passHash []byte, foreignUser bool, foreignUserID string, foreignWebsite proto.Site) (int64, error) {
	// Username is unique, so will fail if user already exists
	var res *sql.Row
	var err error

	switch foreignUser {
	case true:
		res = m.Conn.QueryRow("INSERT INTO users (username, email, pass_hash, foreign_user_ID, foreign_website) "+
			"VALUES ($1, $2, $3, $4, $5) returning id", username, email, string(passHash), foreignUserID, foreignWebsite)
	case false:
		res = m.Conn.QueryRow("INSERT INTO users (username, email, pass_hash) "+
			"VALUES ($1, $2, $3) returning id", username, email, string(passHash))
	}

	var id int64 = 0
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
}

func (m *UserModel) GetUserWithID(userID int64) (*User, error) {
	sql := "SELECT id, username, email FROM users WHERE id=$1"
	user := User{}

	err := m.Conn.Select(&user, sql, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetUserWithUsername(username string) (int64, error) {
	sql := "SELECT id FROM users WHERE username=$1"

	rows := m.Conn.QueryRow(sql, username)

	var userID int64
	err := rows.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *UserModel) GetPassHash(uid int64) (string, error) {
	sql := "SELECT pass_hash FROM users WHERE id = $1"

	row := m.Conn.QueryRow(sql, uid)

	var passHash string
	err := row.Scan(&passHash)
	if err != nil {
		return "", err
	}

	return passHash, nil
}

// Maybe I should cut down on the copy pasta
func (m *UserModel) GetForeignUser(foreignUserID string, foreignWebsite proto.Site) (int64, error) {
	sql := "SELECT id FROM users WHERE foreign_user_ID=$1 AND foreign_website=$2"

	row := m.Conn.QueryRow(sql, foreignUserID, foreignWebsite)

	var userID int64
	err := row.Scan(&userID)

	switch {
	case err == dbsql.ErrNoRows:
		return 0, status.Error(codes.NotFound, errors.UserDoesNotExistMessage)
	case err != nil:
		return 0, fmt.Errorf("scan returned an error: %s", err)
	}

	return userID, nil
}
