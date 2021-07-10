package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	_ "sync"
	"syscall"

	"github.com/horahoradev/horahora/scheduler/internal/models"
	"github.com/horahoradev/horahora/scheduler/internal/syncmanager"

	"github.com/horahoradev/horahora/scheduler/internal/config"
	"github.com/horahoradev/horahora/scheduler/internal/downloader"
	grpcserver "github.com/horahoradev/horahora/scheduler/internal/grpc"
	"github.com/horahoradev/horahora/scheduler/internal/schedule"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Could not get config. Err: %s", err)
	}

	defer cfg.GRPCConn.Close()

	ctx, close := context.WithCancel(context.Background())

	// Handle signals gracefully
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigChan
		log.Errorf("Signal %s received. Canceling context", s)
		close()
	}()

	wg := sync.WaitGroup{}

	dlQueue := make(chan *models.VideoDLRequest, 100)

	// Start one publisher goroutine to poll postgres and send download requests into the channel
	// could potentially expand this to multiple publishers
	wg.Add(1)
	poller, err := schedule.NewPoller(cfg.Conn, cfg.Redlock)
	if err != nil {
		log.Fatalf("Could not create poller. Err: %s", err)
	}

	log.Info("Starting poller")
	go func() {
		err := poller.PollDatabaseAndSendIntoQueue(ctx, dlQueue)
		if err != nil {
			log.Errorf("Database polling failed. Err: %s", err)
			// TODO: Might want to cancel context here since the poller returned, otherwise consumers will just wait indefinitely and perform no work
		}
		wg.Done()
	}()

	m := &sync.Mutex{}
	// Start two goroutines to subscribe to channel and download items
	numOfSubscribers := 3
	for i := 0; i < numOfSubscribers; i++ {
		wg.Add(1)
		dler := downloader.New(dlQueue, cfg.VideoOutputLoc, cfg.Client, cfg.NumberOfRetries, cfg.SocksConnStr, cfg.MaxFS)
		go func() {
			err := dler.SubscribeAndDownload(ctx, m)
			if err != nil {
				log.Errorf("Downloader failed. Err: %s", err)
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		repo := models.NewArchiveRequest(cfg.Conn, cfg.Redlock)

		// TODDO: sync worker exit becausse schcema isn't up yet
		worker, err := syncmanager.NewWorker(repo, cfg.SocksConnStr, cfg.SyncPollDelay)
		if err != nil {
			log.Fatalf("Sync worker exited wth err: %s", err)
		}

		err = worker.Sync()
		if err != nil {
			log.Fatalf("Sync worker exited while syncing with err: %s", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := grpcserver.NewGRPCServer(ctx, cfg.Conn, cfg.Redlock, 7777)
		if err != nil {
			log.Error(err)
		}
		log.Info("GRPC server exited")
	}()

	log.Info("Goroutines started, waiting")
	// Wait for all goroutines to return*/
	wg.Wait()
	log.Info("All goroutines have returned. Exiting...")
}
