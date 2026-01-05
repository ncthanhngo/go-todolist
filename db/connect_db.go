package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to PostgreSQL using GORM and .env config
func ConnectDB() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Lỗi load file .env: ", err)
	}

	// Lấy thông tin từ biến môi trường
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD") // sửa DB_PASS thành DB_PASSWORD cho khớp .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Tạo DSN (Data Source Name) cho Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		host, user, pass, name, port, sslmode,
	)

	// Kết nối Postgres bằng GORM với simple protocol để tránh lỗi encoding
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Tránh lỗi binary encoding
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Kết nối database thất bại:", err)
	}

	log.Println(" Connect to PostgreSQL successfully!")
	return db
}
