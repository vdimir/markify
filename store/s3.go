package store

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io"
	"time"
)

type Config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	bucket          string
}

type S3Storage struct {
	client  *minio.Client
	bucket  string
	timeout time.Duration
}

func NewS3Storage(cfg Config) (*S3Storage, error) {
	minioClient, err := minio.New(cfg.endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.accessKeyID, cfg.secretAccessKey, ""),
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create S3 client")
	}
	exists, err := minioClient.BucketExists(context.Background(), cfg.bucket)
	if err != nil {
		return nil, errors.Wrap(err, "can't access to S3 storage")
	}
	if !exists {
		return nil, fmt.Errorf("bucket %q does not exist", cfg.bucket)
	}
	return &S3Storage{minioClient, cfg.bucket, time.Second * 5}, nil
}

func (s3 *S3Storage) SetBlob(key string, reader io.Reader, meta map[string]string, ttl time.Duration) error {
	opts := minio.PutObjectOptions{}
	if ttl > 0 {
		opts.RetainUntilDate = time.Now().Add(ttl)
	}
	if len(meta) > 0 {
		opts.UserMetadata = meta
	}
	ctx, cancel := context.WithTimeout(context.Background(), s3.timeout)
	defer cancel()
	_, err := s3.client.PutObject(ctx, s3.bucket, key, reader, -1, opts)
	return err
}

func (s3 *S3Storage) GetBlob(key string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s3.timeout)
	defer cancel()
	obj, err := s3.client.GetObject(ctx, s3.bucket, key, minio.GetObjectOptions{})
	return obj, err
}

func (s3 *S3Storage) GetMeta(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s3.timeout)
	defer cancel()
	objMeta, err := s3.client.StatObject(ctx, s3.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return objMeta.UserMetadata, nil
}

func (s3 *S3Storage) DeleteBlob(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s3.timeout)
	defer cancel()
	return s3.client.RemoveObject(ctx, s3.bucket, key, minio.RemoveObjectOptions{})
}
