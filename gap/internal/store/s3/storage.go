package s3

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	Bucket string
}

func New(endpoint string, key string, secret string, bucket string) (*Storage, error) {
	slog.Info("creating new storage client", "endpoint", endpoint, "key", key, "bucket", bucket)

	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(key, secret, ""),
		// if there is a colon, assume its not ssl, ie localhost:9000
		Secure: !strings.Contains(endpoint, ":"),
	})
	if err != nil {
		return nil, err
	}

	store := Storage{client: client, Bucket: bucket}

	err = store.ensureBucket(context.TODO(), bucket)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *Storage) ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) PutReader(ctx context.Context, objectName string, body io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.Bucket, objectName, body, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *Storage) Put(ctx context.Context, objectName string, contentBytes []byte) error {
	buffer := bytes.NewBuffer(contentBytes)
	contentType := mimetype.Detect(contentBytes).String()

	return s.PutReader(ctx, objectName, buffer, int64(len(buffer.Bytes())), contentType)
}

func (s *Storage) GetMinio(ctx context.Context, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, s.Bucket, objectName, minio.GetObjectOptions{})
}

func (s *Storage) Get(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := s.GetMinio(ctx, objectName)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(obj)
}
