package syncmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/horahoradev/horahora/scheduler/internal/models"
	proto "github.com/horahoradev/horahora/scheduler/protocol"
	log "github.com/sirupsen/logrus"
)

type SyncWorker struct {
	R            *models.ArchiveRequestRepo
	SocksConnStr string
	SyncDelay    time.Duration
}

func NewWorker(r *models.ArchiveRequestRepo, socksConnStr string, syncDelay time.Duration) (*SyncWorker, error) {
	return &SyncWorker{R: r,
		SocksConnStr: socksConnStr,
		SyncDelay:    syncDelay}, nil
}

func (s *SyncWorker) Sync() error {
	for {

		dlReqs, err := s.R.GetUnsyncedCategoryDLRequests()
		if err != nil {
			return err
		}

		for _, dlReq := range dlReqs {
			// Distributed lock goes here!

			// refresh cache if backoff period is up
			log.Infof("Backoff period expired for download request %s, syncing all", dlReq.Id)
			itemsAdded, err := s.syncDownloadList(&dlReq)
			if err != nil {
				log.Errorf("Sync worker dl list: %s", err)
				continue
			}

			if itemsAdded {
				err = dlReq.ReportSyncHit()
				if err != nil {
					log.Errorf("sync worker report sync hit: %s", err)
					continue
				}
			} else {
				err = dlReq.ReportSyncMiss()
				if err != nil {
					log.Errorf("Sync worker report sync miss: %s", err)
					continue
				}
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
	URL   string `json:"url"`
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
	cmd := exec.Command("/usr/bin/python3", args...)
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

	log.Errorf("Videos, %+v", videos)

	return videos, nil
}

func (s *SyncWorker) getVideoListString(dlReq *models.CategoryDLRequest) ([]string, error) {
	// TODO: type safety, switch to enum?
	args := []string{"/scheduler/youtube-dl/youtube_dl/__main__.py",
		"-j",
		"--flat-playlist",
	}
	if s.SocksConnStr != "" {
		args = append(args, []string{"--proxy", s.SocksConnStr}...)
	}

	downloadPreference := "all"

	// If it's a tag we're downloading from, then there may be a large number of videos.
	// If we've downloaded from this tag before, we should terminate the search once reaching the latest
	// video we've downloaded.

	// WOW that's a lot of switch statements, should probably flatten or refactor this out into separate functions so
	// that I can actually read this
	switch dlReq.Website {
	case proto.SupportedSite_niconico:
		switch dlReq.ContentType {
		case models.Tag:
			latestVideo, err := dlReq.GetLatestVideoForRequest()

			switch {
			case err == models.NeverDownloaded:
				// keep as all			// TODO: caching download list, lol

				log.Infof("Tag category %s has never been downloaded, downloading all", dlReq.ContentValue)

			case err != nil:
				return nil, err
			default:
				log.Infof("Tag category %s has been downloaded before, resuming at %s", dlReq.ContentValue, *latestVideo)
				downloadPreference = fmt.Sprintf("id%s", *latestVideo)
			}
			args = append(args, fmt.Sprintf("nicosearch%s:%s", downloadPreference, dlReq.ContentValue))

		default:
			err := fmt.Errorf("content type %s is not implemented for niconico.", dlReq.ContentType)
			return nil, err
		}

	case proto.SupportedSite_bilibili:
		switch dlReq.ContentType {
		case models.Tag:
			args = append(args, fmt.Sprintf("bilisearch%s:%s", downloadPreference, dlReq.ContentValue))
			log.Infof("Downloading videos of tag %s from bilibili", dlReq.ContentValue)
			// TODO: implement continuation from latest video for bilibili in extractor
			// for now, try to download everything in the list every time

		case models.Channel:
			log.Infof("Downloading videos from bilibili user %s", dlReq.ContentValue)
			args = append(args, fmt.Sprintf("https://space.bilibili.com/%s", dlReq.ContentValue))

		default:
			err := fmt.Errorf("content type %s is not implemented for bilibili.", dlReq.ContentType)
			return nil, err
		}

	case proto.SupportedSite_youtube:
		// bit of a FIXME: switch this when we've merged both repos into one
		args[0] = "/scheduler/yt-dlp/yt_dlp/__main__.py"

		switch dlReq.ContentType {
		case models.Tag:
			args = append(args, fmt.Sprintf("ytsearch%s:%s", downloadPreference, dlReq.ContentValue))
			log.Infof("downloading videos of tag %s from youtube", dlReq.ContentValue)
		// TODO: ensure youtube extractor returns list in desc order, implements continuation from latest video id

		case models.Channel:
			log.Infof("Downloading videos from youtube user %s", dlReq.ContentValue)
			args = append(args, fmt.Sprintf("https://www.youtube.com/channel/%s", dlReq.ContentValue))

		default:
			err := fmt.Errorf("content type %s is not implemented for youtube.", dlReq.ContentType)
			return nil, err
		}

	default:
		err := fmt.Errorf("no archive request implementations for website %s", dlReq.Website)
		return nil, err
	}

	return args, nil
}
