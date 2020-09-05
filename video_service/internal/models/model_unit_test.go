package models

import (
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLGeneration(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 1, 0, "")
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithUser(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 1, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" WHERE (\"userid\" = 1) ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithTag(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 1, 0, "wow")
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" NATURAL JOIN \"tags\" WHERE (\"tag\" = 'wow') ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}
