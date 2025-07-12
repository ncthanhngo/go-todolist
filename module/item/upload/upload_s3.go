package upload

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"time"
)

func UploadFile(c *gin.Context) {
	// Lấy file từ form-data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc file từ form-data"})
		return
	}
	defer file.Close()

	// Tải cấu hình AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải cấu hình AWS"})
		return
	}

	// Tạo client S3
	s3Client := s3.NewFromConfig(cfg)

	// Lấy thông tin tên file, phần mở rộng
	originalName := header.Filename
	ext := filepath.Ext(originalName)
	uuidName := uuid.New().String()

	// Tạo key dạng: uploads/YYYY/MM/DD/uuid.ext
	datePath := time.Now().Format("2006/01/02")
	key := fmt.Sprintf("uploads/%s/%s%s", datePath, uuidName, ext)

	// Tên bucket
	bucket := "mysoju" // Thay bằng bucket thật của bạn

	// Upload file
	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Upload thất bại: %v", err)})
		return
	}

	// Tạo signed URL có hiệu lực 15 phút
	presigner := s3.NewPresignClient(s3Client)
	presignResult, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải lên thành công nhưng tạo signed URL thất bại"})
		return
	}

	// Trả về kết quả
	c.JSON(http.StatusOK, gin.H{
		"message":    "Tải lên thành công!",
		"file_name":  originalName,
		"stored_key": key,
		"url":        presignResult.URL,
		"expires_in": "15 phút",
	})
}
