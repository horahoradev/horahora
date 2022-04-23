package storage

import "os"

type Storage interface {
	Fetch(path string) (*os.File, error)
	Upload(path, desiredFilename string) error
	Delete(filename string) error
}
