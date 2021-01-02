package downloader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/horahoradev/horahora/scheduler/internal/models"
	proto "github.com/horahoradev/horahora/scheduler/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type downloader struct {
	downloadQueue       chan *models.VideoDlRequest
	outputLoc           string
	videoClient         videoproto.VideoServiceClient
	numberOfRetries     int
	successfulDownloads chan models.Video
}

func New(dlQueue chan *models.VideoDlRequest, outputLoc string, client videoproto.VideoServiceClient, numberOfRetries int, successfulDownloads chan models.Video) downloader {
	return downloader{
		downloadQueue:       dlQueue,
		outputLoc:           outputLoc,
		videoClient:         client,
		numberOfRetries:     numberOfRetries,
		successfulDownloads: successfulDownloads,
	}
}

// SubscribeAndDownload reads from the download queue
// FIXME: should provide slightly better graceful shutdown behavior here
func (d *downloader) SubscribeAndDownload(ctx context.Context) error {
	// This is pretty awkward, but guarantees a return on the next iteration,
	// and guarantees that the function will return if it was waiting on a download
	// request.
	for {
		select {
		case <-ctx.Done():
			log.Info("Context done, downloader returning")
			return nil
		default:
		}

		select {
		case <-ctx.Done():
			log.Info("Context done, downloader returning")
			return nil
		case r := <-d.downloadQueue:
			err := d.downloadRequest(ctx, r)
			if err != nil {
				log.Errorf("Encountered error in downloader. Err: %s. Continuing...", err)
				// FIXME: increase robustness
				//return err
			}
		}
	}
}

type VideoJSON struct {
	Type  string `json:"_type"`
	URL   string `json:"url"`
	IeKey string `json:"ie_key"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Video struct {
	ID string
}

// Deals with a particular download request
func (d *downloader) downloadVideos(ctx context.Context, videos chan Video) error {

	// At this point we have the list of videos to download
	// Iterate over in reverse (ascending by upload date)
	for video := range videos {
		// So this is outrageously dumb, but I'll remove it later. Downloader keeps getting stuck on these bad videos!
		// FIXME
		if strings.HasPrefix(video.ID, "so") {
			log.Info("Video ID has the bad prefix so, skipping...")
			continue
		}
	currVideoLoop:
		for currentRetryNum := 1; currentRetryNum <= d.numberOfRetries+1; currentRetryNum++ {
			select {
			case <-ctx.Done():
				log.Infof("Context done, returning from download request loop for content type %s content val %s", dlReq.ContentType, dlReq.ContentValue)
				return nil
			default:
			}

			switch {
			case currentRetryNum == d.numberOfRetries+1:
				log.Errorf("Failed to download %s in %d attempts", video.URL, d.numberOfRetries)
				break currVideoLoop
			case currentRetryNum > 1:
				log.Infof("Attempting to download %s, attempt %d of %d", video.URL, currentRetryNum, d.numberOfRetries)
			}

			videoReq := videoproto.ForeignVideoCheck{
				ForeignVideoID: video.ID,
				ForeignWebsite: ToVideoSite(dlReq.Website), // LMAO FIXME
			}

			videoExists, err := d.videoClient.ForeignVideoExists(context.TODO(), &videoReq)
			if err != nil {
				err := fmt.Errorf("could not check whether video exists for video ID %s. Err: %s", video.ID, err)
				log.Error(err)
				break currVideoLoop
			}

			if videoExists.Exists {
				log.Errorf("Video ID %s already exists", video.ID)
				break currVideoLoop
			}

			// VideoJSON does not yet exist, try to acquire lock
			err = dlReq.AcquireLockForVideo(video.ID)
			if err != nil {
				// If we can't get the lock, just skip the video in the current archive request
				log.Errorf("Could not acquire redis lock for video ID %s during download of content type %s value %s, err: %s", video.ID,
					dlReq.ContentType, dlReq.ContentValue, err)
				break currVideoLoop
			}

			metafile, metadata, err := d.downloadVideo(video)
			if err == nil {
				log.Infof("Download succeeded for video %s", video.ID)

				// Background is used here to try to ensure that the service will deal with whatever it's currently
				// downloading before shutting down.
				err = d.uploadToVideoService(context.Background(), metadata, video, dlReq.Website, metafile)
				if err != nil {
					log.Infof("failed to upload to video service. Err: %s. Continuing...", err)
					continue
				}

				err := dlReq.SetDownloaded(video.ID)

				// TODO: handle better? retry?
				if err != nil {
					log.Errorf("Could not set latest video. Err: %s", err)
				}

				// TODO: why did I make this channel? I don't remember what it was for, successful video download notifications to client?
				if d.successfulDownloads != nil {
					d.successfulDownloads <- video
				}
				break
			}
			// Just keep trying to download until we succeed
			// TODO: check for specific errors indicating we should skip to the next entry
			log.Errorf("Failed to download video %s. Err: %s", video.ID, err)
		}
	}
	return nil
}

func (d *downloader) syncDownloadList(dlReq *models.VideoDlRequest) (bool, error) {
	videos, err := d.getDownloadList(dlReq)
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

func (d *downloader) getDownloadList(dlReq *models.VideoDlRequest) ([]VideoJSON, error) {
	args, err := getVideoListString(dlReq)
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

	return videos, nil
}

func (d *downloader) downloadVideo(video models.Video) (*os.File, *YTDLMetadata, error) {
	log.Infof("Downloading %s", video)

	args, err := d.getVideoDownloadArgs(&video)
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command("/usr/bin/python3", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Command %s failed with %s. Output: %s", cmd, err, string(output))
		return nil, nil, err
	}

	// surely no video will have a metadata file in excess of 100mb??
	buf := make([]byte, 100*1024*1024)
	file, err := os.Open(fmt.Sprintf("%s/%s.info.json", d.outputLoc, video.ID))
	if err != nil {
		return nil, nil, err
	}

	n, err := file.Read(buf)
	if err != nil {
		return nil, nil, err
	}

	// Truncate
	buf = buf[:n]

	metadata := YTDLMetadata{}

	err = json.Unmarshal(buf, &metadata)
	if err != nil {
		return nil, nil, err
	}

	return file, &metadata, nil
}

// FIXME: this function is quite long and complicated
func (d *downloader) uploadToVideoService(ctx context.Context, metadata *YTDLMetadata, video models.Video,
	website proto.SupportedSite, metafile *os.File) error {
	stream, err := d.videoClient.UploadVideo(ctx)
	if err != nil {
		return fmt.Errorf("could not start video upload stream. Err: %s", err)
	}

	// FIXME: extend to multiple file extensions, this is bad
	generatedVideoFiles, err := filepath.Glob(fmt.Sprintf("%s/%s.mp4", d.outputLoc, video.ID))
	if err != nil {
		return err
	}

	if len(generatedVideoFiles) != 1 {
		return fmt.Errorf("unexpected number of matched files: %d", len(generatedVideoFiles))
	}

	// TODO: fix for other extensions?? this is dumb
	generatedThumbnailFiles, err := filepath.Glob(fmt.Sprintf("%s/%s.jpg", d.outputLoc, video.ID))
	if err != nil {
		return err
	}

	if len(generatedThumbnailFiles) != 1 {
		return fmt.Errorf("unexpected number of matched thumbnail files: %d", len(generatedThumbnailFiles))
	}

	// FIXME: this is dumb
	site := ToVideoSite(website)

	thumb, err := os.Open(generatedThumbnailFiles[0])
	if err != nil {
		return fmt.Errorf("could not open thumbnail. Err: %s", err)
	}
	defer thumb.Close()

	thumbnailContents, err := ioutil.ReadAll(thumb)
	if err != nil {
		return err
	}

	// Send metadata
	metaPayload := videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Meta{
			Meta: &videoproto.InputFileMetadata{
				Title:             metadata.Title,
				Description:       metadata.Description,
				AuthorUID:         metadata.UploaderID,
				OriginalVideoLink: video.URL,
				AuthorUsername:    metadata.Uploader,
				OriginalSite:      site,
				OriginalID:        metadata.ID,
				Tags:              metadata.Tags,
				Thumbnail:         thumbnailContents, // nothing to see here...
			},
		},
	}

	err = stream.Send(&metaPayload)
	if err != nil {
		return fmt.Errorf("could not send metadata. Err: %s", err)
	}

	file, err := os.Open(generatedVideoFiles[0])
	if err != nil {
		return fmt.Errorf("could not open globbed file. Err: %s", err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
		metafile.Close()
		os.Remove(metafile.Name())
	}()

	err = sendLoop(file, stream, false)
	if err != nil {
		return fmt.Errorf("Failed to send video data. Err: %s", err)
	}

	err = sendLoop(metafile, stream, true)
	if err != nil {
		return fmt.Errorf("Failed to send video raw metadata. Err: %s", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("received error after closing stream: %s", err)
	}

	log.Infof("Video %s has been uploaded as video ID %d", video.URL, resp.VideoID)
	return nil
}

func sendLoop(file *os.File, stream videoproto.VideoService_UploadVideoClient, isMeta bool) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

loop:
	for {
		buf := make([]byte, 1*1024*1024)
		n, err := file.Read(buf)

		switch {
		// I think it's fine to check for EOF and no n==0, but just in case...
		case n == 0 && err == io.EOF:
			break loop
		case err != io.EOF && err != nil:
			return fmt.Errorf("could not read from file. Err: %s", err)
		}

		// Truncate
		buf = buf[:n]

		var chunkPayload videoproto.InputVideoChunk
		switch isMeta {
		case true:
			chunkPayload = videoproto.InputVideoChunk{
				Payload: &videoproto.InputVideoChunk_Rawmeta{
					Rawmeta: &videoproto.RawMetadata{
						Data: buf,
					},
				},
			}

		case false:
			chunkPayload = videoproto.InputVideoChunk{
				Payload: &videoproto.InputVideoChunk_Content{
					Content: &videoproto.FileContent{
						Data: buf,
					},
				},
			}

		}

		err = stream.Send(&chunkPayload)
		switch {
		case err == io.EOF:
			return fmt.Errorf("videoservice closed stream prematurely")
		case err != nil:
			return fmt.Errorf("could not send to videoservice. Err: %s", err)
		}
	}

	return nil
}

func getVideoListString(dlReq *models.VideoDlRequest) ([]string, error) {
	// TODO: type safety, switch to enum?
	args := []string{"/scheduler/youtube-dl/youtube_dl/__main__.py", "-j", "--flat-playlist"}
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
				// keep as all
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

func (d *downloader) getVideoDownloadArgs(video *models.Video) ([]string, error) {
	args := []string{
		"/scheduler/youtube-dl/youtube_dl/__main__.py",
		video.URL,
		"--write-info-json", // I'd like to use -j, but doesn't seem to work for some videos
		"--write-thumbnail",
		"--get-comments",
		"-o",
		fmt.Sprintf("%s/%s", d.outputLoc, "%(id)s.%(ext)s"),
	}

	return args, nil
}

type YTDLMetadata struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Formats []struct {
		URL         string `json:"url"`
		Ext         string `json:"ext"`
		FormatID    string `json:"format_id"`
		FormatNote  string `json:"format_note,omitempty"`
		Container   string `json:"container,omitempty"`
		Quality     int    `json:"quality"`
		Filesize    int    `json:"filesize,omitempty"`
		Format      string `json:"format"`
		Protocol    string `json:"protocol"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			AcceptCharset  string `json:"Accept-Charset"`
			Accept         string `json:"Accept"`
			AcceptEncoding string `json:"Accept-Encoding"`
			AcceptLanguage string `json:"Accept-Language"`
			Cookie         string `json:"Cookie"`
		} `json:"http_headers"`
		Abr    float64 `json:"abr,omitempty"`
		Vbr    float64 `json:"vbr,omitempty"`
		Height int     `json:"height,omitempty"`
		Width  int     `json:"width,omitempty"`
		Tbr    float64 `json:"tbr,omitempty"`
	} `json:"formats"`
	Thumbnails []struct {
		URL string `json:"url"`
		Ext string `json:"ext"`
		ID  string `json:"id"`
	} `json:"thumbnails"`
	Description  string   `json:"description"`
	Uploader     string   `json:"uploader"`
	Timestamp    int      `json:"timestamp"`
	UploaderID   string   `json:"uploader_id"`
	ViewCount    int      `json:"view_count"`
	Tags         []string `json:"tags"`
	Genre        string   `json:"genre"`
	CommentCount int      `json:"comment_count"`
	RawComments  struct {
		En []struct {
			Ping struct {
				Content string `json:"content"`
			} `json:"ping,omitempty"`
			Thread struct {
				Resultcode    int    `json:"resultcode"`
				Thread        string `json:"thread"`
				ServerTime    int    `json:"server_time"`
				LastRes       int    `json:"last_res"`
				Ticket        string `json:"ticket"`
				Revision      int    `json:"revision"`
				ClickRevision int    `json:"click_revision"`
			} `json:"thread,omitempty"`
			Leaf struct {
				Thread string `json:"thread"`
				Count  int    `json:"count"`
			} `json:"leaf,omitempty"`
			Chat struct {
				Thread    string `json:"thread"`
				Language  int    `json:"language"`
				No        int    `json:"no"`
				Vpos      int    `json:"vpos"`
				Date      int    `json:"date"`
				Premium   int    `json:"premium"`
				Anonymity int    `json:"anonymity"`
				UserID    string `json:"user_id"`
				Mail      string `json:"mail"`
				Content   string `json:"content"`
			} `json:"chat,omitempty"`
		} `json:"en"`
		Jp []struct {
			Ping struct {
				Content string `json:"content"`
			} `json:"ping,omitempty"`
			Thread struct {
				Resultcode    int    `json:"resultcode"`
				Thread        string `json:"thread"`
				ServerTime    int    `json:"server_time"`
				LastRes       int    `json:"last_res"`
				Ticket        string `json:"ticket"`
				Revision      int    `json:"revision"`
				ClickRevision int    `json:"click_revision"`
			} `json:"thread,omitempty"`
			Leaf struct {
				Thread string `json:"thread"`
				Count  int    `json:"count"`
			} `json:"leaf,omitempty"`
			Chat struct {
				Thread    string `json:"thread"`
				No        int    `json:"no"`
				Vpos      int    `json:"vpos"`
				Leaf      int    `json:"leaf"`
				Date      int    `json:"date"`
				Anonymity int    `json:"anonymity"`
				UserID    string `json:"user_id"`
				Mail      string `json:"mail"`
				Content   string `json:"content"`
			} `json:"chat,omitempty"`
		} `json:"jp"`
		Cn []struct {
			Ping struct {
				Content string `json:"content"`
			} `json:"ping,omitempty"`
			Thread struct {
				Resultcode    int    `json:"resultcode"`
				Thread        string `json:"thread"`
				ServerTime    int    `json:"server_time"`
				LastRes       int    `json:"last_res"`
				Ticket        string `json:"ticket"`
				Revision      int    `json:"revision"`
				ClickRevision int    `json:"click_revision"`
			} `json:"thread,omitempty"`
			Leaf struct {
				Thread string `json:"thread"`
				Count  int    `json:"count"`
			} `json:"leaf,omitempty"`
			Chat struct {
				Thread    string `json:"thread"`
				Language  int    `json:"language"`
				No        int    `json:"no"`
				Vpos      int    `json:"vpos"`
				Leaf      int    `json:"leaf"`
				Date      int    `json:"date"`
				Anonymity int    `json:"anonymity"`
				UserID    string `json:"user_id"`
				Mail      string `json:"mail"`
				Content   string `json:"content"`
			} `json:"chat,omitempty"`
		} `json:"cn"`
	} `json:"raw_comments"`
	Comments []struct {
		Parent    interface{} `json:"parent"`
		ID        int         `json:"id"`
		AuthorID  string      `json:"author_id"`
		Text      string      `json:"text"`
		Timestamp int         `json:"timestamp"`
		Language  string      `json:"language"`
	} `json:"comments"`
	Subtitles struct {
		DanmakuEn []struct {
			Ext  string `json:"ext"`
			Data string `json:"data"`
		} `json:"danmaku-en"`
		DanmakuJp []struct {
			Ext  string `json:"ext"`
			Data string `json:"data"`
		} `json:"danmaku-jp"`
		DanmakuCn []struct {
			Ext  string `json:"ext"`
			Data string `json:"data"`
		} `json:"danmaku-cn"`
	} `json:"subtitles"`
	Duration           float64     `json:"duration"`
	WebpageURL         string      `json:"webpage_url"`
	Extractor          string      `json:"extractor"`
	WebpageURLBasename string      `json:"webpage_url_basename"`
	ExtractorKey       string      `json:"extractor_key"`
	Playlist           interface{} `json:"playlist"`
	PlaylistIndex      interface{} `json:"playlist_index"`
	Thumbnail          string      `json:"thumbnail"`
	DisplayID          string      `json:"display_id"`
	UploadDate         string      `json:"upload_date"`
	URL                string      `json:"url"`
	FormatID           string      `json:"format_id"`
	Ext                string      `json:"ext"`
	Abr                float64     `json:"abr"`
	Vbr                float64     `json:"vbr"`
	Height             int         `json:"height"`
	Width              int         `json:"width"`
	Quality            int         `json:"quality"`
	Tbr                float64     `json:"tbr"`
	Format             string      `json:"format"`
	Protocol           string      `json:"protocol"`
	HTTPHeaders        struct {
		UserAgent      string `json:"User-Agent"`
		AcceptCharset  string `json:"Accept-Charset"`
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
	} `json:"http_headers"`
	Fulltitle string `json:"fulltitle"`
	Filename  string `json:"_filename"`
}

// TODO: remove in future
func ToVideoSite(website proto.SupportedSite) videoproto.Website {
	return videoproto.Website(videoproto.Website_value[website.String()])
}
