package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/horahoradev/horahora/backup_service/internal/config"
	"github.com/horahoradev/horahora/video_service/storage"
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

	for true {
		backupAllDBs(conf, storageAPI)
		time.Sleep(time.Hour)
	}
}

func backupAllDBs(c *config.Config, storageAPI *storage.B2Storage) error {
	err := dumpAndWrite("userservice", c.UserPGSDatabase, c.UserPGSUsername, c.UserPGSPassword, storageAPI)
	if err != nil {
		return err
	}

	err = dumpAndWrite("scheduler", c.SchedulerPGSDatabase, c.SchedulerPGSUsername, c.SchedulerPGSPassword, storageAPI)
	if err != nil {
		return err
	}

	err = dumpAndWrite("videoservice", c.VideoPGSDatabase, c.VideoPGSUsername, c.VideoPGSPassword, storageAPI)
	if err != nil {
		return err
	}
	return nil
}

func dumpAndWrite(hostname, dbName, username, password string, b2 *storage.B2Storage) error {
	filename := fmt.Sprintf("%s-%d", dbName, time.Now().Unix())

	dumpFile, err := os.Create(fmt.Sprintf("/tmp/%s", filename))
	if err != nil {
		return err
	}
	defer dumpFile.Close()

	cmd := exec.Command("/usr/bin/pg_dump", []string{
		fmt.Sprintf("--dbname=postgresql://%s:%s@%s:5432/%s", username, password, hostname, dbName),
	}...)

	cmd.Stdout = dumpFile

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = b2.Upload(dumpFile.Name(), filename)
	if err != nil {
		return err
	}
	return nil
}
