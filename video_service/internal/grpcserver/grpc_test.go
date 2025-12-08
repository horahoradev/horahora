package grpcserver

import (
	"io"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jmoiron/sqlx"

	"github.com/golang/mock/gomock"

	userproto "github.com/horahoradev/horahora/user_service/protocol"
	usermocks "github.com/horahoradev/horahora/user_service/protocol/mocks"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	mocks "github.com/horahoradev/horahora/video_service/protocol/mocks"
	"github.com/stretchr/testify/assert"
)

// This is a really huge test... it'd probably be better if I split things up
func TestGRPCUpload(t *testing.T) {
	bucketName := "horahora-dev-otomads"

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxMock := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	// https://github.com/DATA-DOG/go-sqlmock/issues/27
	mock.ExpectQuery("INSERT INTO videos").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	mock.ExpectCommit()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResp := userproto.GetForeignUserResponse{NewUID: 1}

	mockClient := usermocks.NewMockUserServiceClient(mockCtrl)
	mockClient.EXPECT().GetUserForForeignUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockResp, nil)

	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	g, err := initGRPCServer(bucketName, sqlxMock, mockClient, true, "http://localhost", "s3", "", "", 0, "", 0, rdb, 50)
	assert.NoError(t, err)

	file, err := os.Create("../../test_files/NO.mp4")
	assert.NoError(t, err)
	filedata := make([]byte, 1024*1024*5) // 5 megs is enough
	_, err = file.Write(filedata)
	assert.NoError(t, err)
	_, err = file.Seek(0, 0)
	assert.NoError(t, err)

	_, err = file.Read(filedata)
	assert.NoError(t, err)

	mockServ := mocks.NewMockVideoService_UploadVideoServer(mockCtrl)
	payload := videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Content{
			Content: &videoproto.FileContent{
				Data: filedata,
			},
		},
	}

	metaPayload := videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Meta{
			Meta: &videoproto.InputFileMetadata{
				Title:             "test_title",
				Description:       "wow",
				AuthorUID:         "1",
				AuthorUsername:    "test_user",
				OriginalVideoLink: "www.nicovideo.jp/watch/sm9",
				OriginalID:        "sm9",
				OriginalSite:      "niconico",
			},
		},
	}

	mockServ.EXPECT().Recv().Return(&metaPayload, nil).Times(1)
	mockServ.EXPECT().Recv().Return(&payload, nil).Times(1)
	mockServ.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	mockServ.EXPECT().SendAndClose(gomock.Any()).Return(nil).Times(1)

	err = g.UploadVideo(mockServ)
	assert.NoError(t, err)
}
