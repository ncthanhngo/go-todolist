package ginitem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"todolist/common"
	"todolist/module/item/biz"
	"todolist/module/item/storage"
)

func GetItemById(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {

		//Doan nay thu nghiem Recovery - khong lien quan toi handler
		go func() {
			defer common.Recovery() // Chan panic o runtime de chuong trinh tiep tuc chayj
			var a []int
			fmt.Println(a[0]) // -> error panic
		}()

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		store := storage.NewSQLStore(db)
		business := biz.NewGetItemBiz(store)

		data, err := business.GetItem(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
