package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// ConnectDB connects to MySQL using GORM and .env config
func ConnectDB() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(" Lỗi load file .env: ", err)
	}

	// Lấy thông tin từ biến môi trường
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Tạo DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	// Kết nối MySQL bằng GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Kết nối database thất bại:", err)
	}

	log.Println("✅ Connect to database successfully!")
	return db
}
