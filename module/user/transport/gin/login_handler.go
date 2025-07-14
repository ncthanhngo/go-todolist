package ginuser

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todolist/common"
	"todolist/component/tokenprovider/jwt"
	biz2 "todolist/module/user/biz"
	"todolist/module/user/model"
	"todolist/module/user/storage"
)

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUerData model.UserLogin
		if err := c.ShouldBind(&loginUerData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		tokenProvider := jwt.NewTokenJWTProvider("jwt", "ncthanh@999") //truyen bien moi truong vao
		store := storage.NewSQLStore(db)
		bcypt := common.NewBcryptHasher(10)

		business := biz2.NewLoginBusiness(store, tokenProvider, bcypt, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUerData)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
