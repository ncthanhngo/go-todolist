package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"todolist/db"
	"todolist/middleware"
	"todolist/module/upload/storage"
	"todolist/module/upload/transport"

	//"gorm.io/driver/mysql"
	//"gorm.io/gorm"
	"log"
	ginitem "todolist/module/item/transport/gin"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("DB_USER = ", os.Getenv("DB_USER")) // Debug

}
func main() {

	db := db.ConnectDB()
	db = db.Debug()

	// Chay Gin Framework
	r := gin.Default()
	r.Use(middleware.Recover()) //middleware se tac dong toi toan bo API ben duoi
	r.Static("/static", "./static")
	v1 := r.Group("v1")
	//s3
	//r.PUT("/upload", upload.UploadFile)
	//s3 Clean
	// Cấu hình AWS S3

	r.PUT("/upload", transport.UploadFileHandler)
	//Khoi {} phia duoi la khoi tu do, gioi han va tao su de nhin
	//Khai bao dang ky cho 5 API
	{

		//v1.PUT("/upload", upload.Upload(db))
		items := v1.Group("/items")
		{
			items.POST("", ginitem.CreateItem(db)) // ginitem la sua package trong handler tranh trung voi gin
			items.GET("", ginitem.ListItem(db))
			items.GET("/:id", ginitem.GetItemById(db))
			items.PATCH("/:id", ginitem.UpdateItemByID(db))
			items.DELETE("/:id", ginitem.DeleteItemByID(db))
		}
	}
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

// Week vs Strong: Bang Strong la bang co nhieu khoa ngoai tham chieu toi no
//Week  thuong la cac bang giua ket noi cac moi quan he n-n
