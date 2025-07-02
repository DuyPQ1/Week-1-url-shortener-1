package main

import (
	"url-shortener/database"
	"url-shortener/handlers"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Kết nối tới database
	db := database.ConnectDB()

	// Tự động tạo bảng theo model URLMapping
	db.AutoMigrate(&models.URLMapping{})

	// Khởi tạo Gin router
	r := gin.Default()

	// Route để rút gọn URL
	r.POST("/shorten", handlers.ShortenURLHandler(db))

	// Route để chuyển hướng từ URL rút gọn về URL gốc
	r.GET("/:code", handlers.RedirectHandler(db))

	// Route để xem thống kê số lượt click
	r.GET("/stats/:code", handlers.StatsHandler(db))

	// Chạy server trên port 2000
	r.Run(":2000")
}
