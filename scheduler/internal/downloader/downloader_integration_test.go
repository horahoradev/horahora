package downloader

import (
	"context"
	"github.com/horahoradev/horahora/scheduler/internal/config"
	"github.com/horahoradev/horahora/scheduler/internal/models"
	proto "github.com/horahoradev/horahora/scheduler/protocol"
	log "github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

// Tests whether we can succeed at downloading something, upload to videoservice, and return when context is cancelled
func TestConcurrentVideoDownloads(t *testing.T) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	dlChan := make(chan *models.CategoryDLRequest, 10)

	dlSuccessChannel := make(chan VideoJSON, 5)
	downloader := New(dlChan, "./", conf.Client, 5, dlSuccessChannel)

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	go func() {
		wg.Add(1)
		err := downloader.SubscribeAndDownload(ctx)
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	dlChan <- models.NewVideoDlRequest(proto.SupportedSite_niconico, models.Tag, "Oldest_Video", "1", conf.Conn, conf.Redlock)

	timeoutTicker := time.Tick(30 * time.Minute)

	select {
	case <-timeoutTicker:
		log.Fatal("wait for successful download timed out")

	case <-dlSuccessChannel:
	}

	// Cancel context, causing downloader to return after it's done with the current video download
	cancel()

	wg.Wait()
	<-dlSuccessChannel
}
