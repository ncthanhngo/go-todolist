package transport

import (
	"net/http"
	"todolist/module/upload/storage"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc file"})
			return
		}
		defer file.Close()
		uploader, err := storage.NewS3Uploader("ap-southeast-1", "mysoju")
		if err != nil {
			panic(err)
		}
		key, url, err := uploader.UploadFile(c.Request.Context(), file, header.Filename, header.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Tải lên thành công",
			"key":     key,
			"url":     url,
		})
	}
}
