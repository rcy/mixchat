package space

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Space struct {
	s3Client    *s3.S3
	bucket      string
	uriEndpoint string
}

type InitParams struct {
	S3Key       string
	S3Secret    string
	Endpoint    string
	URIEndpoint string
	Bucket      string
}

func MustInit(params InitParams) *Space {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(params.S3Key, params.S3Secret, ""),
		Endpoint:         aws.String(params.Endpoint),
		Region:           aws.String("us-east-1"), // not actually for DO, but need to supply this here
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		panic(err)
	}

	space := Space{
		s3Client:    s3.New(newSession),
		bucket:      params.Bucket,
		uriEndpoint: params.URIEndpoint,
	}

	fmt.Printf("Initialized s3 storage %s, %s, %s\n", space.s3Client.Endpoint, space.bucket, space.uriEndpoint)

	return &space
}

func (s *Space) Put(ctx context.Context, key string, data []byte) error {
	//log.Printf("Space.Put %s", key)

	object := s3.PutObjectInput{
		ContentType: aws.String(http.DetectContentType(data)),
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ACL:         aws.String("public-read"),
		Metadata:    map[string]*string{
			//"x-amz-meta-my-key": aws.String("your-value"),
		},
	}

	_, err := s.s3Client.PutObjectWithContext(ctx, &object)
	return err
}

func (s *Space) Get(ctx context.Context, key string) ([]byte, error) {
	//log.Printf("Space.Get %s", key)

	object := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	out, err := s.s3Client.GetObjectWithContext(ctx, &object)
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (s *Space) URI(key string) string {
	return s.uriEndpoint + "/" + key
}
