package grpcserver

import (
	"io"
	"os"
	"testing"

	userproto "github.com/horahoradev/horahora/video_service/protocol"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jmoiron/sqlx"

	"github.com/golang/mock/gomock"

	usermocks "github.com/horahoradev/horahora/user_service/protocol/mocks"
	proto "github.com/horahoradev/horahora/video_service/protocol"
	mocks "github.com/horahoradev/horahora/video_service/protocol/mocks"
	"github.com/stretchr/testify/assert"

	_ "github.com/stretchr/testify/assert"
)

// This is a really huge test... it'd probably be better if I split things up
func TestGRPCUpload(t *testing.T) {
	bucketName := "horahora-dev-videos"

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxMock := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO videos").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResp := userproto.GetForeignUserResponse{
		NewUID: 1,
	}

	mockClient := usermocks.NewMockUserServiceClient(mockCtrl)
	mockClient.EXPECT().GetUserForForeignUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockResp, nil)

	g, err := initGRPCServer(bucketName, sqlxMock, mockClient)
	assert.NoError(t, err)

	file, err := os.Open("../../test_files/NO.mp4")
	assert.NoError(t, err)

	filedata := make([]byte, 1024*1024*5) // 5 megs is enough

	_, err = file.Read(filedata)
	assert.NoError(t, err)

	mockServ := mocks.NewMockVideoService_UploadVideoServer(mockCtrl)
	payload := proto.InputVideoChunk{
		Payload: &proto.InputVideoChunk_Content{
			Content: &proto.FileContent{
				Data: filedata,
			},
		},
	}

	metaPayload := proto.InputVideoChunk{
		Payload: &proto.InputVideoChunk_Meta{
			Meta: &proto.InputFileMetadata{
				Title:             "test_title",
				Description:       "wow",
				AuthorUID:         "1",
				OriginalVideoLink: "www.nicovideo.jp/watch/sm9",
				OriginalSite:      proto.Website_niconico,
			},
		},
	}

	mockServ.EXPECT().Recv().Return(&metaPayload, nil).Times(1)
	mockServ.EXPECT().Recv().Return(&payload, nil).Times(1)
	mockServ.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	err = g.UploadVideo(mockServ)
	assert.NoError(t, err)
}
