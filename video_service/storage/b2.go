package storage

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/kurin/blazer/b2"
)

/*
Copyright 2016, the Blazer authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

I modified some of the source code used in Blazer's README when writing this file.

Summary of changes:
- April 3, 2021: made some small changes to to the source examples to get the code to comply with my storage interface

*/

type B2Storage struct {
	Client *b2.Client
	Bucket *b2.Bucket
}

func NewB2(backblazeID, backblazeAPIKey, bucketname string) (*B2Storage, error) {
	// b2_authorize_account
	client, err := b2.NewClient(context.Background(), backblazeID, backblazeAPIKey)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(context.Background(), bucketname)
	if err != nil {
		return nil, err
	}

	return &B2Storage{Client: client, Bucket: bucket}, nil

}

func (s *B2Storage) Fetch(id string) (*os.File, error) {
	r := s.Bucket.Object(id).NewReader(context.Background())
	defer r.Close()

	f, err := os.Create(uploadDir + id)
	if err != nil {
		return nil, err
	}

	r.ConcurrentDownloads = 1
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}

func (s *B2Storage) Upload(path, desiredFilename string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	obj := s.Bucket.Object(desiredFilename)
	w := obj.NewWriter(context.Background())
	if _, err := io.Copy(w, f); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

// TODO: Not implemented
func (s *B2Storage) Delete(filename string) error {
	return errors.New("Not implemented") // there's no easy delete operation
}
