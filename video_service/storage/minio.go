package storage

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinio(minioEndpoint, minioKeyID, minioAccessKey, bucketname string) (*MinioStorage, error) {
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioKeyID, minioAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{Client: minioClient, Bucket: bucketname}, nil
}

func (s *MinioStorage) Fetch(id string) (*os.File, error) {
	r, err := s.Client.GetObject(context.Background(), s.Bucket, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer r.Close()

	f, err := os.Create(uploadDir + id)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *MinioStorage) Upload(path, desiredFilename string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	_, err = s.Client.PutObject(context.Background(), s.Bucket, desiredFilename, f, fStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
