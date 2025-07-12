package biz

import (
	"context"
	"io"
)

type Uploader interface {
	UploadFile(ctx context.Context, file io.Reader, filename, contentType string) (key string, url string, err error)
}
