package syncmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/horahoradev/horahora/scheduler/internal/models"
	log "github.com/sirupsen/logrus"
)

type SyncWorker struct {
	R                 *models.ArchiveRequestRepo
	SocksConnStr      string
	SyncDelay         time.Duration
	RequestDLCountMap map[string]int
}

func NewWorker(r *models.ArchiveRequestRepo, socksConnStr string, syncDelay time.Duration) (*SyncWorker, error) {
	return &SyncWorker{R: r,
		SocksConnStr:      socksConnStr,
		SyncDelay:         syncDelay,
		RequestDLCountMap: make(map[string]int),
	}, nil
}

func (s *SyncWorker) Sync() error {
	for {

		dlReqs, err := s.R.GetUnsyncedCategoryDLRequests()
		if err != nil {
			return err
		}

		for _, dlReq := range dlReqs {
			// TODO: distributed lock goes here!

			// refresh cache if backoff period is up
			log.Infof("Backoff period expired for download request %s, syncing all", dlReq.Id)
			itemsAdded, err := s.syncDownloadList(&dlReq)

			if itemsAdded {
				err := dlReq.ReportSyncHit()
				if err != nil {
					log.Errorf("sync worker report sync hit: %s", err)
				}
			} else {
				err := dlReq.ReportSyncMiss()
				if err != nil {
					log.Errorf("Sync worker report sync miss: %s", err)
				}
			}

			if err != nil {
				log.Errorf("Sync worker dl list: %s", err)
				continue
			}

		}

		time.Sleep(s.SyncDelay)
	}
}

func (s *SyncWorker) syncDownloadList(dlReq *models.CategoryDLRequest) (bool, error) {
	videos, err := s.getDownloadList(dlReq)
	if err != nil {
		return false, fmt.Errorf("could not fetch download list. Err: %s", err)
	}

	var newItemsAdded bool

	for _, video := range videos {

		// TODO: batch?
		itemsAdded, err := dlReq.AddVideo(video.ID, video.URL)
		if err != nil {
			return false, fmt.Errorf("could not add video. Err: %s", err)
		}

		if itemsAdded {
			newItemsAdded = true
		}
	}

	return newItemsAdded, nil
}

type VideoJSON struct {
	Type  string `json:"_type"`
	URL   string `json:"original_url"`
	IeKey string `json:"ie_key"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *SyncWorker) getDownloadList(dlReq *models.CategoryDLRequest) ([]VideoJSON, error) {
	args, err := s.getVideoListString(dlReq)
	if err != nil {
		return nil, err
	}

	// get the list of videos to download
	cmd := exec.Command(args[0], args[1:]...)
	payload, err := cmd.Output()
	if err != nil {
		log.Errorf("Command `%s` finished with err %s", cmd, err)
		return nil, err
	}

	var videos []VideoJSON
	// The json isn't formatted as a list LMAO please FIXME
	// I assume that the list provided by youtube-dl will be in descending order by upload date.
	// Download by upload date asc so that we can resume at newest download.
	spl := strings.Split(string(payload), "\n")
	for i := len(spl) - 2; i >= 0; i-- {
		line := spl[i]
		var video VideoJSON
		log.Infof("Line: %s", line)
		err = json.Unmarshal([]byte(line), &video)

		if err != nil {
			log.Errorf("Failed to unmarshal json. Payload: %s. Err: %s", line, err)
			return nil, err
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		log.Errorf("Could not unmarshal, videolist len is 0")
		return nil, errors.New("unmarshal failure")
	}

	return videos, nil
}

func (s *SyncWorker) RunVideoClassificationLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Second * 30)
		}

		urls, err := s.R.GetUnclassifiedVideoURLs()
		if err != nil {
			return err
		}

		for _, url := range urls {
			classification, err := s.GetVideoClassification(url.URL)
			if err != nil {
				continue
			}

			err = s.R.UpdateClassification(classification, url.ID)
			if err != nil {
				return err
			}
		}
	}
}

func (s *SyncWorker) GetVideoClassification(videoURL string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("yt-dlp --add-header 'Accept-Language:ja' -j %s | jq '.tags'", videoURL))
	payload, err := cmd.Output()
	if err != nil {
		log.Errorf("Command `%s` finished with err %s", cmd, err)
		return "", err
	}

	categories, err := s.R.GetInferenceCategories()
	if err != nil {
		return "", err
	}

	for _, category := range categories {
		if strings.Index(string(payload), category.Tag) != -1 {
			return category.Category, nil
		}
	}

	return "General", nil
}

func (s *SyncWorker) getVideoListString(dlReq *models.CategoryDLRequest) ([]string, error) {
	// TODO: type safety, switch to enum?
	args := []string{"yt-dlp",
		"-j",
		"--flat-playlist",
	}

	_, ok := s.RequestDLCountMap[dlReq.Id]
	if !ok {
		s.RequestDLCountMap[dlReq.Id] = 0
	}

	if s.RequestDLCountMap[dlReq.Id]%10 != 0 {
		args = append(args, []string{
			"--playlist-end",
			"400",
		}...)
	}

	s.RequestDLCountMap[dlReq.Id]++

	if s.SocksConnStr != "" {
		args = append(args, []string{"--proxy", s.SocksConnStr}...)
	}

	args[0] = "yt-dlp"

	args = append(args, dlReq.Url)
	return args, nil
}
