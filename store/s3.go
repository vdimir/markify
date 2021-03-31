package store

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io"
	"strings"
	"time"
)

type S3Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key"`
	SecretAccessKey string `json:"secret"`
	Bucket          string `json:"bucket"`
}

type S3Storage struct {
	client  *minio.Client
	bucket  string
	timeout time.Duration
	ctx context.Context
}

func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create S3 client")
	}
	exists, err := minioClient.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		return nil, errors.Wrap(err, "can't access to S3 storage")
	}
	if !exists {
		return nil, fmt.Errorf("bucket %q does not exist", cfg.Bucket)
	}
	return &S3Storage{minioClient, cfg.Bucket, time.Second * 5, context.Background()}, nil
}

func (s3 *S3Storage) SetBlob(key string, reader io.Reader, meta map[string]string, ttl time.Duration) error {
	opts := minio.PutObjectOptions{}
	if ttl > 0 {
		opts.RetainUntilDate = time.Now().Add(ttl)
	}
	if len(meta) > 0 {
		opts.UserMetadata = meta
	}
	_, err := s3.client.PutObject(s3.ctx, s3.bucket, key, reader, -1, opts)
	return errors.Wrap(err, "s3 put object error")
}

func (s3 *S3Storage) GetBlob(key string) (io.Reader, map[string]string, error) {
	obj, err := s3.client.GetObject(s3.ctx, s3.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "s3 get object error")
	}
	stat, err := obj.Stat()
	if err != nil {
		return nil, nil, errors.Wrap(err, "s3 get object meta error")
	}
	mapKeysToLower(stat.UserMetadata)
	return obj, stat.UserMetadata, nil
}

func (s3 *S3Storage) GetMeta(key string) (map[string]string, error) {
	objMeta, err := s3.client.StatObject(s3.ctx, s3.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "s3 metadata error")
	}
	mapKeysToLower(objMeta.UserMetadata)
	return objMeta.UserMetadata, nil
}

func (s3 *S3Storage) DeleteBlob(key string) error {
	return s3.client.RemoveObject(s3.ctx, s3.bucket, key, minio.RemoveObjectOptions{})
}

func mapKeysToLower(m map[string]string) {
	for k, v := range m {
		m[strings.ToLower(k)] = v
	}
}