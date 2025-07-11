package ginitem

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	"todolist/module/item/biz"
	"todolist/module/item/model"
	"todolist/module/item/storage"
)

func ListItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}
		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		queryString.Paging.Process()

		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store)
		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
