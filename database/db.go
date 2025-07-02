package database

import (
	"log"
	"url-shortener/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Hàm kết nối tới MySQL database
func ConnectDB() *gorm.DB {
	// Chuỗi kết nối MySQL với thông tin: user:password@protocol(host:port)/database
	dsn := "root:12345@tcp(127.0.0.1:3306)/url_database?charset=utf8mb4&parseTime=True&loc=Local"

	// Mở kết nối tới database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("DB Error:", err)
	}

	// Tự động migrate bảng URLMapping
	if err := db.AutoMigrate(&models.URLMapping{}); err != nil {
		// Nếu có lỗi trong quá trình migrate, dừng chương trình và in thông báo lỗi
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
