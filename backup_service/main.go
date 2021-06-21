package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/horahoradev/horahora/backup_service/internal/config"
	"github.com/horahoradev/horahora/video_service/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	storageAPI, err := storage.NewB2(conf.BackblazeID, conf.BackblazeAPIKey, "otomads")
	if err != nil {
		log.Fatal(err)
	}

	// Crude synchronization mechanism
	time.Sleep(time.Second * 15)

	for true {
		err = backupAllDBs(conf, storageAPI)
		if err != nil {
			log.Errorf("Failed to backup databases. Err: %s", err)
		}
		time.Sleep(time.Hour)
	}
}

func backupAllDBs(c *config.Config, storageAPI *storage.B2Storage) error {
	log.Print("Starting to dump databases...")
	err := dumpAndWrite("postgres", c.UserPGSDatabase, c.UserPGSUsername, c.UserPGSPassword, storageAPI)
	if err != nil {
		return err
	}
	log.Print("Dumped userservice")

	err = dumpAndWrite("postgres", c.SchedulerPGSDatabase, c.SchedulerPGSUsername, c.SchedulerPGSPassword, storageAPI)
	if err != nil {
		return err
	}
	log.Print("Dumped scheduler")

	err = dumpAndWrite("postgres", c.VideoPGSDatabase, c.VideoPGSUsername, c.VideoPGSPassword, storageAPI)
	if err != nil {
		return err
	}
	log.Print("Dumped videoservice")

	log.Print("Database dump complete.")
	return nil
}

func dumpAndWrite(hostname, dbName, username, password string, b2 *storage.B2Storage) error {
	filename := fmt.Sprintf("backup_%s-%d.sql", dbName, time.Now().Unix())

	dumpFile, err := os.Create(fmt.Sprintf("/tmp/%s", filename))
	if err != nil {
		return err
	}

	defer func() {
		// Fine to remove the file locally after it's been uploaded to the storage backend
		os.Remove(dumpFile.Name())
		dumpFile.Close()
	}()

	log.Printf("Dumping to %s", filename)

	cmd := exec.Command("/usr/bin/pg_dump", []string{
		fmt.Sprintf("--dbname=postgresql://%s:%s@%s:5432/%s", username, password, hostname, dbName),
	}...)

	cmd.Stdout = dumpFile

	err = cmd.Run()
	if err != nil {
		log.Errorf("Command %s failed with %s.", cmd, err)
		return err
	}

	log.Printf("Uploading %s to b2", dumpFile.Name())

	err = b2.Upload(dumpFile.Name(), filename)
	if err != nil {
		return err
	}
	return nil
}
