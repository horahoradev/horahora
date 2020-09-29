package models

import (
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLGeneration(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "", true)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithUser(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 2, 1, "", true)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" WHERE (\"userid\" = 1) ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithTag(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "wow", true)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" INNER JOIN \"tags\" ON (\"videos\".\"id\" = \"tags\".\"video_id\") WHERE (\"tag\" = 'wow') ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithTagApprovedOnly(t *testing.T) {
	sql, err := generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "wow", false)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\" FROM \"videos\" INNER JOIN \"tags\" ON (\"videos\".\"id\" = \"tags\".\"video_id\") WHERE ((\"tag\" = 'wow') AND (\"is_approved\" IS TRUE)) ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}
