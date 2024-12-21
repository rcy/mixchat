package s3

import (
	"bytes"
	"context"
	"testing"
)

func TestUpload(t *testing.T) {
	ctx := context.Background()
	bucketName := "bucket1234"
	objectName := "object5678"
	contentBytes := []byte("abc def")

	s, err := New("localhost:9500", "minioadmin", "minioadmin", bucketName)
	if err != nil {
		t.Fatalf("error creating storage: %s", err)
	}

	err = s.ensureBucket(ctx, bucketName)
	if err != nil {
		t.Fatalf("error making bucket: %s", err)
	}

	err = s.Put(ctx, objectName, contentBytes)
	if err != nil {
		t.Fatalf("error uploading: %s", err)
	}

	downloadedBytes, err := s.Get(ctx, objectName)
	if err != nil {
		t.Fatalf("error downloading: %s", err)
	}

	if !bytes.Equal(contentBytes, downloadedBytes) {
		t.Fatalf("downloaded does not match original")
	}
}
