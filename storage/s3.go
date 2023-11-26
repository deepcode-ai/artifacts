package storage

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ClientOpts struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	UseSSL          bool   `json:"useSSL"`
}

type S3StorageClient struct {
	minioClient *minio.Client
}

// NewS3StorageClient unmarshals the S3 storage credentials and then initializes a new S3StorageClient.
func NewS3StorageClient(_ context.Context, credentialsJSON []byte) (*S3StorageClient, error) {
	s3StorageCredentials := &S3ClientOpts{}
	if err := json.Unmarshal(credentialsJSON, s3StorageCredentials); err != nil {
		return nil, err
	}

	minioClient, err := minio.New(s3StorageCredentials.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3StorageCredentials.AccessKeyID, s3StorageCredentials.SecretAccessKey, ""),
		Secure: s3StorageCredentials.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3StorageClient{minioClient}, nil
}

// UploadDir uploads a directory to the specified bucket.
func (s *S3StorageClient) UploadDir(bucket, src, dst string) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = s.UploadDir(bucket, filepath.Join(src, file.Name()), filepath.Join(dst, file.Name()))
			if err != nil {
				return err
			}
		} else {
			err = s.UploadObject(bucket, filepath.Join(src, file.Name()), filepath.Join(dst, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UploadObject uploads a single object to the specified bucket.
func (s *S3StorageClient) UploadObject(bucket, src, dst string) (err error) {
	_, err = s.minioClient.FPutObject(context.Background(), bucket, dst, src, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GetObjects downloads a list of objects from the specified bucket.
func (s *S3StorageClient) GetObjects(bucket, destinationPath string, paths ...string) error {
	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(path string) {
			defer wg.Done()

			err := s.GetObject(bucket, path, filepath.Join(destinationPath, path))
			if err != nil {
				return
			}
		}(path)
	}
	return nil
}

// GetObject downloads a single object from the specified bucket.
func (s *S3StorageClient) GetObject(bucket, src, dst string) (err error) {
	err = s.minioClient.FGetObject(context.Background(), bucket, src, dst, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// NewReader returns a new io.ReadCloser for the specified object.
func (s *S3StorageClient) NewReader(ctx context.Context, bucket, src string) (io.ReadCloser, error) {
	return s.minioClient.GetObject(ctx, bucket, src, minio.GetObjectOptions{})
}
