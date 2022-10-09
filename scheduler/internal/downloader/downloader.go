package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/horahoradev/horahora/scheduler/internal/models"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	log "github.com/sirupsen/logrus"
)

type downloader struct {
	downloadQueue   chan *models.VideoDLRequest
	outputLoc       string
	videoClient     videoproto.VideoServiceClient
	numberOfRetries int
	socksConnStr    string
	maxFS           uint64
	acceptLanguage  string
}

func New(dlQueue chan *models.VideoDLRequest, outputLoc string, client videoproto.VideoServiceClient, numberOfRetries int,
	socksConnStr string, maxFS uint64, acceptLanguage string) downloader {
	return downloader{
		downloadQueue:   dlQueue,
		outputLoc:       outputLoc,
		videoClient:     client,
		numberOfRetries: numberOfRetries,
		socksConnStr:    socksConnStr,
		maxFS:           maxFS,
		acceptLanguage:  acceptLanguage,
	}
}

// SubscribeAndDownload reads from the download queue
// FIXME: should provide slightly better graceful shutdown behavior here
func (d *downloader) SubscribeAndDownload(ctx context.Context, m *sync.Mutex) error {
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
			err := d.downloadVideoReq(ctx, r, m)
			if err != nil {
				log.Errorf("Encountered error in downloader. Err: %s. Continuing...", err)
				// FIXME: increase robustness
				//return err
			}
		}
	}
}

// Deals with a particular video download request
func (d *downloader) downloadVideoReq(ctx context.Context, video *models.VideoDLRequest, m *sync.Mutex) error {
	if strings.HasPrefix(video.VideoID, "so") {
		err := video.SetDownloadFailed()
		if err != nil {
			log.Errorf("Could not set download failed for video %s. Err: %s", video.VideoID, err)
		}
		log.Info("Video VideoID has the bad prefix so, skipping...")
		return nil
	}

	err := video.SetDownloadInProgress()
	if err != nil {
		log.Errorf("Failed to set download in progress: %v", err)
	}

	website, err := models.GetWebsiteFromURL(video.URL)
	if err != nil {
		log.Errorf("Failed to extract website domain from %s", video.URL)
		return nil
	}

	videoReq := videoproto.ForeignVideoCheck{
		ForeignVideoID: video.VideoID,
		ForeignWebsite: website, // lol FIXME
	}

	videoExists, err := d.videoClient.ForeignVideoExists(context.TODO(), &videoReq)
	if err != nil {
		err := fmt.Errorf("could not check whether video exists for video VideoID %s. Err: %s", video.VideoID, err)
		log.Error(err)
		return nil
	}

	if videoExists.Exists {
		log.Errorf("Video VideoID %s already exists", video.VideoID)
		err = video.SetDownloadSucceeded()
		if err != nil {
			log.Errorf("Could not set download succeeded for video %s. Err: %s", video.VideoID, err)
		}
		return nil
	}

	// LOL
	// there's a race condition in which ffprobe can stall indefinitely if its cookies are invalidated before its stream metadata
	// request succeeds. This is my attempt to lower the likelihood of this issue occurring.
	m.Lock()
	go func() {
		time.Sleep(time.Second * 10)
		m.Unlock()
	}()

	// The + 10 is just a precaution against channel blockages in case i've overlooked something
	errCh := make(chan error, d.numberOfRetries+10)
currVideoLoop:
	for currentRetryNum := 1; currentRetryNum <= d.numberOfRetries+1; currentRetryNum++ {
		select {
		case <-ctx.Done():
			log.Infof("Context done, returning from download request loop for parent url", video.ParentURL)
			return nil
		default:
		}

		switch {
		case currentRetryNum == d.numberOfRetries+1:
			log.Errorf("Failed to download %s in %d attempts", video.URL, d.numberOfRetries)
			err := video.SetDownloadFailed()
			if err != nil {
				log.Errorf("Could not set download failed for video %s. Err: %s", video.VideoID, err)
			}

			close(errCh)
			err = <-errCh // just get the last error for now, TODO
			err = video.RecordEvent(models.Error, err.Error())
			if err != nil {
				log.Errorf("Could not record error event. Err: %s", err)
			}

			break currVideoLoop
		case currentRetryNum > 1:
			log.Infof("Attempting to download %s, attempt %d of %d", video.URL, currentRetryNum, d.numberOfRetries)
		}

		metafile, metadata, err := d.downloadVideo(video)
		if err == nil {
			log.Infof("Download succeeded for video %s", video.VideoID)

			// Background is used here to try to ensure that the service will deal with whatever it's currently
			// downloading before shutting down.
			err = d.uploadToVideoService(context.Background(), metadata, video, metafile)
			if err != nil {
				errCh <- err
				log.Infof("failed to upload to video service. Err: %s. Continuing...", err)
				continue
			}

			err = video.SetDownloadSucceeded()
			if err != nil {
				errCh <- err
				log.Errorf("Could not set download succeeded for video %s. Err: %s", video.VideoID, err)
			}

			//err := video.SetDownloaded()
			//
			//// TODO: handle better? retry?
			//if err != nil {
			//	log.Errorf("Could not set latest video. Err: %s", err)
			//}

			break
		}
		// Just keep trying to download until we succeed
		// TODO: check for specific errors indicating we should skip to the next entry
		errCh <- err
		log.Errorf("Failed to download video %s. Err: %s", video.VideoID, err)
	}
	return nil
}

func (d *downloader) downloadVideo(video *models.VideoDLRequest) (*os.File, *YTDLMetadata, error) {
	log.Infof("Downloading %v+", video)

	args, err := d.getVideoDownloadArgs(video)
	if err != nil {
		return nil, nil, err
	}

	ytdlLog, err := os.Create(fmt.Sprintf("%s/%s.ytdl", d.outputLoc, video.VideoID))
	if err != nil {
		return nil, nil, err
	}
	defer ytdlLog.Close()

	ctxTimeout, _ := context.WithTimeout(context.Background(), time.Second*900)

	cmd := exec.CommandContext(ctxTimeout, args[0], args[1:]...)
	cmd.Stdout = ytdlLog
	cmd.Stderr = ytdlLog

	err = cmd.Run()
	if err != nil {
		log.Errorf("Command %s failed with %s.", cmd, err)
		return nil, nil, err
	}

	file, err := os.Open(fmt.Sprintf("%s/%s.info.json", d.outputLoc, video.VideoID))
	if err != nil {
		return nil, nil, err
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	metadata := YTDLMetadata{}

	err = json.Unmarshal(buf, &metadata)
	if err != nil {
		return nil, nil, err
	}

	return file, &metadata, nil
}

// FIXME: this function is quite long and complicated
func (d *downloader) uploadToVideoService(ctx context.Context, metadata *YTDLMetadata, video *models.VideoDLRequest, metafile *os.File) error {
	stream, err := d.videoClient.UploadVideo(ctx)
	if err != nil {
		return fmt.Errorf("could not start video upload stream. Err: %s", err)
	}

	var (
		videoExts = []string{"mp4", "webm", "flv", "mkv"}
		thumbExts = []string{"png", "webp", "jpg"}
	)

	// FIXME: extend to multiple file extensions, this is bad, not DRY
	// I could use walk but I don't want to. I will fix this eventually!
	generatedVideoFiles, err := globWithExtensions(fmt.Sprintf("%s/*%s", d.outputLoc, video.VideoID), videoExts)
	if err != nil {
		return err
	}

	// TODO: fix for other extensions?? this is dumb
	generatedThumbnailFiles, err := globWithExtensions(fmt.Sprintf("%s/*%s", d.outputLoc, video.VideoID), thumbExts)
	if err != nil {
		return err
	}

	thumb, err := os.Open(generatedThumbnailFiles[0])
	if err != nil {
		return fmt.Errorf("could not open thumbnail. Err: %s", err)
	}
	defer thumb.Close()

	thumbnailContents, err := ioutil.ReadAll(thumb)
	if err != nil {
		return err
	}

	// TODO: we need more robust screening for null essential fields
	// https://github.com/envoyproxy/protoc-gen-validate looks promising!
	if metadata.UploaderID == "" {
		metadata.UploaderID = metadata.ChannelID
	}

	website, err := models.GetWebsiteFromURL(video.URL)
	if err != nil {
		log.Errorf("failed to extract website from url: %s", video.URL)
	}

	// Send metadata
	// REFACTOR TODO
	metaPayload := videoproto.InputVideoChunk{
		Payload: &videoproto.InputVideoChunk_Meta{
			Meta: &videoproto.InputFileMetadata{
				Title:             metadata.Title,
				Description:       metadata.Description,
				AuthorUID:         metadata.UploaderID,
				OriginalVideoLink: video.URL,
				AuthorUsername:    metadata.Uploader,
				OriginalSite:      website,
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
		return fmt.Errorf("failed to send video data. Err: %s", err)
	}

	err = sendLoop(metafile, stream, true)
	if err != nil {
		return fmt.Errorf("failed to send video raw metadata. Err: %s", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("received error after closing stream: %s", err)
	}

	log.Infof("Video %s has been uploaded as video VideoID %d", video.URL, resp.VideoID)
	err = video.RecordEvent(models.Downloaded, "")
	if err != nil {
		log.Errorf("Could not record downloaded event. Err: %s. Continuing...", err)
	}

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

func (d *downloader) getVideoDownloadArgs(video *models.VideoDLRequest) ([]string, error) {
	acceptLanguageString := "Accept-Language:en"
	if d.acceptLanguage != "" {
		acceptLanguageString = fmt.Sprintf("Accept-Language:%s", d.acceptLanguage)
	}
	bin := "yt-dlp"
	args := []string{
		bin,
		video.URL,
		"--write-info-json", // I'd like to use -j, but doesn't seem to work for some videos
		"--write-thumbnail",
		// This line was originally authored by Soichiro
		// according to him, it was licensed under the
		// "Do What The Fuck You Want license", for which
		// usage is permitted as long as the name is changed
		// Thank you for your work!
		"-S",
		"res,hdr,fps,vcodec:av01:h265:vp9.2:vp9:h264,vbr",
		"--add-header",
		"Accept:*/*",
		// "Why do we need this?"
		// Previously ffprobe would stall indefinitely if nico's cookies were invalidated by the time it made a request
		// (or something like that).
		"--add-header",
		acceptLanguageString,
		"--socket-timeout",
		"1800",
		"--verbose",
		"-o",
		// Some websites have two IDs per video, so I made it explicit just to avoid issues
		fmt.Sprintf("%s/%s.%s", d.outputLoc, video.VideoID, "%(ext)s"),
	}

	if d.maxFS != 0 {
		args = append(args, []string{"--max-filesize", fmt.Sprintf("%dm", d.maxFS)}...)
	}

	// FIXME: This is dumb
	// youtube doesn't support get-comments
	if !strings.HasPrefix(video.URL, "https://www.youtube.com") {
		args = append(args, "--get-comments")

	}

	if d.socksConnStr != "" {
		args = append(args, []string{"--proxy", d.socksConnStr}...)
	}

	return args, nil
}

func globWithExtensions(basepath string, extensions []string) ([]string, error) {
	var ret []string
	for _, ext := range extensions {
		path := fmt.Sprintf("%s.%s*", basepath, ext)
		g, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		ret = append(ret, g...)
	}

	if len(ret) != 1 {
		return nil, fmt.Errorf("incorrect match length for extensions %s, length: %d, contents: %s", extensions, len(ret), ret)
	}

	return ret, nil
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
	ChannelID    string   `json:"channel_id"`
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
	//Comments []struct {
	//	Parent    interface{} `json:"parent"`
	//	ID        int         `json:"id"`
	//	AuthorID  string      `json:"author_id"`
	//	Text      string      `json:"text"`
	//	Timestamp int         `json:"timestamp"`
	//	Language  string      `json:"language"`
	//} `json:"comments"`
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
