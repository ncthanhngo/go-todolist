package ginitem

import (
	"net/http"
	"strconv"
	"todolist/common"
	"todolist/module/item/biz"
	"todolist/module/item/model"
	"todolist/module/item/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateItemByID(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var updateData model.TodoItemUpdate

		if err := c.ShouldBind(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		store := storage.NewSQLStore(db)
		business := biz.NewUpdateItemBiz(store, requester)
		if err := business.UpdateItemById(c.Request.Context(), id, &updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
