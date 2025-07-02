package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Hàm tạo mã ngắn ngẫu nhiên với độ dài n ký tự
func generateShortCode(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, n)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// Handler xử lý việc rút gọn URL
func ShortenURLHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Struct để nhận JSON request
		var req struct {
			URL string `json:"url"`
		}

		// Kiểm tra dữ liệu đầu vào
		if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Kiểm tra xem URL đã tồn tại chưa
		var existing models.URLMapping
		if err := db.Where("long_url = ?", req.URL).First(&existing).Error; err == nil {
			// Nếu đã tồn tại, trả về URL ngắn cũ
			c.JSON(http.StatusOK, gin.H{
				"short_url": fmt.Sprintf("http://localhost:2000/%s", existing.ShortCode),
				"note":      "URL đã tồn tại, dùng short URL cũ",
			})
			return
		}

		// Tạo mã ngắn mới
		shortCode := generateShortCode(6)
		entry := models.URLMapping{
			LongURL:   req.URL,
			ShortCode: shortCode,
			CreatedAt: time.Now(),
			Clicks:    0,
		}

		// Lưu vào database
		if err := db.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		// Trả về URL đã rút gọn
		c.JSON(http.StatusOK, gin.H{
			"short_url": fmt.Sprintf("http://localhost:2000/%s", shortCode),
		})
	}
}

// Handler xử lý việc chuyển hướng từ URL ngắn về URL gốc
func RedirectHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy mã ngắn từ URL parameter
		code := c.Param("code")

		// Tìm URL mapping trong database
		var entry models.URLMapping
		if err := db.Where("short_code = ?", code).First(&entry).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Tăng số lượt click
		db.Model(&entry).UpdateColumn("clicks", entry.Clicks+1)

		// Chuyển hướng về URL gốc
		c.Redirect(http.StatusFound, entry.LongURL)
	}
}

// Handler xử lý việc xem thống kê số lượt click
func StatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy mã ngắn từ URL parameter
		code := c.Param("code")

		// Tìm URL mapping trong database
		var entry models.URLMapping
		if err := db.Where("short_code = ?", code).First(&entry).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Trả về thông tin thống kê
		c.JSON(http.StatusOK, gin.H{
			"url":    entry.LongURL,
			"clicks": entry.Clicks,
		})
	}
}
