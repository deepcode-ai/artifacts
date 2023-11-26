package storage

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleCloudStorageClient struct {
	client *storage.Client
}

func NewGoogleCloudStorageClient(ctx context.Context, credentialsJSON []byte) (*GoogleCloudStorageClient, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, err
	}

	return &GoogleCloudStorageClient{client}, nil
}

func (s *GoogleCloudStorageClient) UploadDir(bucket, src, dst string) error {
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

func (s *GoogleCloudStorageClient) UploadObject(bucket, src, dst string) (err error) {
	file, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	obj := s.client.Bucket(bucket).Object(dst)
	w := obj.NewWriter(context.Background())
	if _, err = w.Write(file); err != nil {
		log.Printf("error uploading file %q: %v", dst, err)
		return
	}
	if err = w.Close(); err != nil {
		log.Printf("error closing writer for file %q: %v", dst, err)
		return
	}
	return
}

func (s *GoogleCloudStorageClient) GetObjects(bucket string, destinationPath string, paths ...string) error {
	for _, path := range paths {
		obj := s.client.Bucket(bucket).Object(path)
		r, err := obj.NewReader(context.Background())
		if err != nil {
			return err
		}

		defer r.Close()

		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		err = os.WriteFile(destinationPath, data, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GoogleCloudStorageClient) NewReader(ctx context.Context, bucket string, path string) (io.ReadCloser, error) {
	obj := s.client.Bucket(bucket).Object(path)
	return obj.NewReader(ctx)
}
