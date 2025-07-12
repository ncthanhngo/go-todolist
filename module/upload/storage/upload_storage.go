package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"path/filepath"
	"strings"
	"time"
	"todolist/module/upload/biz"
)

type s3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Uploader(region, bucket string) (biz.Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)

	return &s3Uploader{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

func (s *s3Uploader) UploadFile(ctx context.Context, file io.Reader, filename, contentType string) (string, string, error) {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	datePath := time.Now().Format("2006/01/02")
	key := fmt.Sprintf("uploads/%s/%s-%d%s", datePath, base, time.Now().UnixNano(), ext)

	uploader := manager.NewUploader(s.client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", fmt.Errorf("upload thất bại: %v", err)
	}

	// URL tĩnh nếu bucket public
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)

	return key, url, nil
}
