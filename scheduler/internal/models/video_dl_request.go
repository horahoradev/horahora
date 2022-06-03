package models

import (
	"encoding/json"
	"fmt"
	"time"

	stomp "github.com/go-stomp/stomp/v3"
	log "github.com/sirupsen/logrus"

	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
)

type VideoDLRequest struct {
	Redsync     *redsync.Redsync
	Rabbitmq    *stomp.Conn
	Db          *sqlx.DB
	VideoID     string // Foreign ID
	ID          int    // Domestic ID
	DownloaddID int
	URL         string
	ParentURL   string
	mut         *redsync.Mutex
}

func (v *VideoDLRequest) SetDownloadSucceeded() error {
	sql := "UPDATE videos SET dlStatus = 1 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	if err != nil {
		return err
	}
	// Publish into rabbitmq
	return v.PublishVideoInprogress(1, "deletion")
}

func (v *VideoDLRequest) SetDownloadFailed() error {
	sql := "UPDATE videos SET dlStatus = 2 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	if err != nil {
		return err
	}
	// Publish into rabbitmq
	return v.PublishVideoInprogress(2, "deletion")
}

const queueName = "/topic/videosinprogress"

func (v *VideoDLRequest) SetDownloadInProgress() error {
	sql := "UPDATE videos SET dlStatus = 3 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	if err != nil {
		return err
	}
	// Publish into rabbitmq
	return v.PublishVideoInprogress(3, "insertion")
}

// have to pass a transaction for this one because it needs to be atomic with the scheduler query
func (v *VideoDLRequest) SetDownloadQueued() error {
	sql := "UPDATE videos SET dlStatus = 4 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	if err != nil {
		return err
	}
	// Publish into rabbitmq
	return v.PublishVideoInprogress(4, "insertion")
}

type VideoProgress struct {
	VideoID  string
	Website  string
	DlStatus string
}

type ProgressNotification struct {
	Type  string
	Video VideoProgress
}

func (v *VideoDLRequest) PublishVideoInprogress(dlStatus int, action string) error {
	website, err := GetWebsiteFromURL(v.ParentURL)
	if err != nil {
		return err
	}

	dlStatusString := ""
	if dlStatus == 3 {
		dlStatusString = "Downloading"
	} else if dlStatus == 4 {
		dlStatusString = "Queued"
	}

	p := ProgressNotification{
		Video: VideoProgress{
			VideoID:  v.VideoID,
			Website:  website,
			DlStatus: dlStatusString,
		},
		Type: action, // insertion or deletion
	}

	log.Infof("Publishing %s", p)

	payload, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	return v.Rabbitmq.Send(queueName, "text/plain", payload, nil)
}

func (v *VideoDLRequest) AcquireLockForVideo() error {
	v.mut = v.Redsync.NewMutex(v.VideoID, redsync.SetExpiry(time.Minute*10))
	return v.mut.Lock()
}

func (v *VideoDLRequest) ReleaseLockForVideo() error {
	if v.mut != nil {
		_, err := v.mut.Unlock()
		return err
	}
	return nil

}

type event string

const (
	Scheduled  event = "Video %s from %s has been scheduled for download"
	Error      event = "Video %s from %s could not be downloaded, failed with an error. "
	Downloaded event = "Video %s from %s has been downloaded successfully, and uploaded to videoservice"
)

func (v *VideoDLRequest) RecordEvent(inpEvent event, additionalErrorMsg string) error {
	website, err := GetWebsiteFromURL(v.ParentURL)
	if err != nil {
		return err
	}

	formattedMsg := fmt.Sprintf(string(inpEvent), v.VideoID, website)

	if additionalErrorMsg != "" {
		formattedMsg += fmt.Sprintf("\n\nError message: %s", additionalErrorMsg)
	}

	sql := "insert into archival_events (video_url, download_id, parent_url, event_message, event_time) VALUES ($1, $2, $3, $4, Now())"
	_, err = v.Db.Exec(sql, v.URL, v.DownloaddID, v.ParentURL, formattedMsg)
	return err
}
