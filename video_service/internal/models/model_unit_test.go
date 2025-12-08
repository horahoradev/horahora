package models

import (
	"testing"

	"github.com/golang/mock/gomock"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	usermocks "github.com/horahoradev/horahora/user_service/protocol/mocks"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSQLGeneration(t *testing.T) {
	v := &VideoModel{}
	sql, err := v.generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "", true, false, videoproto.OrderCategory_upload_date)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\", \"views\" FROM \"videos\" WHERE ((\"transcoded\" IS TRUE) AND (\"is_deleted\" IS FALSE)) ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithUser(t *testing.T) {
	v := &VideoModel{}
	sql, err := v.generateVideoListSQL(videoproto.SortDirection_asc, 2, 1, "", true, false, videoproto.OrderCategory_upload_date)
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT \"videos\".\"id\", \"title\", \"userid\", \"newlink\", \"views\" FROM \"videos\" WHERE ((\"userid\" = 1) AND (\"transcoded\" IS TRUE) AND (\"is_deleted\" IS FALSE)) ORDER BY \"upload_date\" ASC LIMIT 50 OFFSET 50")
}

func TestSQLGenerationWithTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := usermocks.NewMockUserServiceClient(ctrl)

	v := &VideoModel{
		grpcClient: mockClient,
	}

	// Mock the call that happens inside getConditions for "username" search term check
	// actually getConditions iterates over cols: title, tag, username.
	// For "username", it calls GetUserIDsForUsername.
	// We need to expect that call.
	// Wait, the search term is "wow". It checks all 3 columns.
	// So it WILL call GetUserIDsForUsername("wow%").

	mockClient.EXPECT().GetUserIDsForUsername(gomock.Any(), &userproto.GetUserIDsForUsernameRequest{Username: "wow%"}).Return(&userproto.GetUserIDsForUsernameResponse{UserIDs: []int64{1}}, nil)

	sql, err := v.generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "wow", true, false, videoproto.OrderCategory_upload_date)
	assert.NoError(t, err)
	assert.NotEmpty(t, sql)
}

func TestSQLGenerationWithTagApprovedOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := usermocks.NewMockUserServiceClient(ctrl)

	v := &VideoModel{
		grpcClient: mockClient,
	}

	mockClient.EXPECT().GetUserIDsForUsername(gomock.Any(), &userproto.GetUserIDsForUsernameRequest{Username: "wow%"}).Return(&userproto.GetUserIDsForUsernameResponse{UserIDs: []int64{1}}, nil)

	sql, err := v.generateVideoListSQL(videoproto.SortDirection_asc, 2, 0, "wow", false, false, videoproto.OrderCategory_upload_date)
	assert.NoError(t, err)
	assert.NotEmpty(t, sql)
}
