package ginitem

import (
	"net/http"
	"todolist/common"
	"todolist/module/item/biz"
	"todolist/module/item/model"
	"todolist/module/item/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.TodoItemCreation

		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		itemData.UserId = requester.GetUserId()
		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}
