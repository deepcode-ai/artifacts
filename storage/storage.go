package storage

import (
	"context"
	"fmt"
	"io"
)

type StorageClient interface {
	UploadDir(string, string, string) error
	UploadObject(string, string, string) error
	GetObjects(string, string, ...string) error
	NewReader(context.Context, string, string) (io.ReadCloser, error)
}

func NewStorageClient(ctx context.Context, storageType string, credentials []byte) (StorageClient, error) {
	switch storageType {
	case "gcs":
		return NewGoogleCloudStorageClient(ctx, credentials)
	case "s3":
		return NewS3StorageClient(ctx, credentials)
	default:
		return &GoogleCloudStorageClient{}, fmt.Errorf("expected storageType to be 'gcs' or 's3'. Received %s", storageType)
	}
}
