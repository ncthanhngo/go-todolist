package ginuser

import (
	"net/http"
	"todolist/common"

	"github.com/gin-gonic/gin"
)

func Profile() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
