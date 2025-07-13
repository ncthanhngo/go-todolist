package ginuser

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	biz2 "todolist/module/user/biz"
	"todolist/module/user/model"
	"todolist/module/user/storage"
)

func Register(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var data model.UserCreate
		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}
		store := storage.NewSQLStore(db)
		bcryp := common.NewBcryptHasher(10)
		biz := biz2.NewRegisterBusiness(store, bcryp)
		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}

}
