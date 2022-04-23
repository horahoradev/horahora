package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	log "github.com/sirupsen/logrus"
)

const uploadDir = "/tmp/"

type S3Storage struct {
	BucketName string
	S3Client   s3.Client
}

func NewS3(bucketName string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(cfg)

	return &S3Storage{S3Client: *s3Client, BucketName: bucketName}, nil
}

func NewS3Minio(bucketName string, endpoint string) (*S3Storage, error) {
	// Credit goes to danmux on Stackoverflow for this specific snippet:
	// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
	// Thank you for your work!

	const defaultRegion = "us-east-1"
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: defaultRegion,
		}, nil
	})

	cfg := aws.Config{
		Region:           defaultRegion,
		Credentials:      credentials.NewStaticCredentialsProvider("minioadmin", "minioadmin", ""),
		EndpointResolver: staticResolver,
	}

	s3Client := s3.New(cfg)
	return &S3Storage{S3Client: *s3Client, BucketName: bucketName}, nil
}

// FIXME: can probably rewrite a significant portion of this. Too long and complicated!
func (s *S3Storage) Fetch(id string) (*os.File, error) {
	getReq := &s3.GetObjectInput{
		Bucket: &s.BucketName,
		Key:    &id,
	}

	getObjReq := s.S3Client.GetObjectRequest(getReq)
	res, err := getObjReq.Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := os.OpenFile(uploadDir+id, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

loop:
	for true {
		buf := make([]byte, 1024*1024*1)
		n, err := res.Body.Read(buf)

		switch {
		case n == 0 && err == io.EOF:
			break loop
		case err != io.EOF && err != nil:
			err = fmt.Errorf("could not read from S3 for transcoding. Err: %s", err)
			log.Error(err)
			f.Close()
			return nil, err
		}

		// Truncate
		buf = buf[:n]

		_, err = f.Write(buf)
		if err != nil {
			f.Close()
			return nil, err
		}
	}

	return f, nil
}

func (s *S3Storage) Upload(path, desiredFilename string) error {
	data, err := os.Open(path)
	if err != nil {
		return err
	}

	putObjInp := s3.PutObjectInput{
		ACL:                       s3.ObjectCannedACLPublicRead,
		Body:                      data,
		Bucket:                    &s.BucketName,
		CacheControl:              nil,
		ContentDisposition:        nil,
		ContentEncoding:           nil,
		ContentLanguage:           nil,
		ContentLength:             nil,
		ContentMD5:                nil,
		ContentType:               nil, // TODO
		Expires:                   nil,
		GrantFullControl:          nil,
		GrantRead:                 nil,
		GrantReadACP:              nil,
		GrantWriteACP:             nil,
		Key:                       &desiredFilename,
		Metadata:                  nil,
		ObjectLockLegalHoldStatus: "",
		ObjectLockMode:            "",
		ObjectLockRetainUntilDate: nil,
		RequestPayer:              "",
		SSECustomerAlgorithm:      nil,
		SSECustomerKey:            nil,
		SSECustomerKeyMD5:         nil,
		SSEKMSEncryptionContext:   nil,
		SSEKMSKeyId:               nil,
		ServerSideEncryption:      "",
		StorageClass:              s3.StorageClassStandardIa,
		Tagging:                   nil,
		WebsiteRedirectLocation:   nil,
	}

	putObjReq := s.S3Client.PutObjectRequest(&putObjInp)
	_, err = putObjReq.Send(context.TODO())
	return err // TODO
}

func (s *S3Storage) Delete(filename string) error {
	deleteObjInp := s3.DeleteObjectInput{
		Bucket: &s.BucketName,
		Key:    &filename,
	}

	delObjReq := s.S3Client.DeleteObjectRequest(&deleteObjInp)
	_, err := delObjReq.Send(context.TODO())
	return err
}
