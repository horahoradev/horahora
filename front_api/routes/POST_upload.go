package routes

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (r RouteHandler) upload(c echo.Context) error {
	profile, err := r.getUserProfileInfo(c)
	if err != nil {
		return err
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	tagsInt := c.FormValue("tags")

	var tags []string
	err = json.Unmarshal([]byte(tagsInt), &tags)
	if err != nil {
		return err
	}

	thumbFileHeader, err := c.FormFile(thumbnailKey)
	if err != nil {
		return err
	}

	thumbFile, err := thumbFileHeader.Open()
	if err != nil {
		return err
	}

	videoFileHeader, err := c.FormFile(videoKey)
	if err != nil {
		return err
	}

	videoFile, err := videoFileHeader.Open()
	if err != nil {
		return err
	}

	// TODO: rewrite me, this isn't memory efficient
	videoBytes, err := ioutil.ReadAll(videoFile)
	if err != nil {
		return err
	}

	thumbBytes, err := ioutil.ReadAll(thumbFile)
	if err != nil {
		return err
	}

	log.Infof("Title: %s", title)
	log.Infof("Description: %s", description)
	log.Infof("Tags: %s", tags)
	log.Infof("Video len: %d thumb len: %d", len(videoBytes), len(thumbBytes))

	uploadClient, err := r.v.UploadVideo(context.Background())
	if err != nil {
		return err
	}

	metaChunk := &videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Meta{
			Meta: &videoproto.InputFileMetadata{
				Title:       title,
				Description: description,
				//AuthorUID:            "",
				//OriginalVideoLink:    "",
				//AuthorUsername:       "",
				//OriginalSite:         0,
				//OriginalID:           "",
				DomesticAuthorID: profile.UserID,
				Tags:             tags,
				Thumbnail:        thumbBytes,
			},
		},
	}

	err = uploadClient.Send(metaChunk)
	if err != nil {
		return err
	}

	for byteInd := 0; byteInd < len(videoBytes); byteInd += fileUploadChunkSize {
		videoByteSlice := videoBytes[byteInd:min(len(videoBytes), byteInd+fileUploadChunkSize)]
		log.Infof("uploading byte %d", byteInd)
		videoChunk := &videoproto.InputVideoChunk{
			Payload: &videoproto.InputVideoChunk_Content{
				Content: &videoproto.FileContent{
					Data: videoByteSlice,
				},
			},
		}

		err = uploadClient.Send(videoChunk)
		if err != nil {
			return err
		}
	}

	resp, err := uploadClient.CloseAndRecv()
	if err != nil {
		return err
	}

	// Redirect to the new video
	return c.JSON(http.StatusOK, resp.VideoID)
}
