package ginitem

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"todolist/common"
	"todolist/module/item/biz"
	"todolist/module/item/storage"
)

func DeleteItemByID(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Xóa thẳng ay (mất tích trong database, chỉ dành cho thưc thể yếu)
		//	if err := db.Table(TodoItem{}.TableName()).Where("id =?", id).Delete(nil).Error; err != nil {
		//Soft Delete se update dong thoi, vi du chuyen status thanh deleted
		store := storage.NewSQLStore(db)
		business := biz.NewDeleteItemBiz(store)
		if err := business.DeleteItem(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
