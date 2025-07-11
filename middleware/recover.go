package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todolist/common"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					return
				}
				// Nếu err không phải AppError, trả về lỗi nội bộ
				appErr := common.ErrInternal(fmt.Errorf("%v", err))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
			}
		}()
		c.Next()
	}
}

// vao Main chay r.Use(middleware.recover())
