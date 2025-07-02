package models

import (
	"time"

	"gorm.io/gorm"
)

// Model đại diện cho bảng lưu trữ mapping giữa URL dài và URL ngắn
type URLMapping struct {
	gorm.Model           // Embed các field cơ bản của GORM (ID, CreatedAt, UpdatedAt, DeletedAt)
	ID         uint      `json:"id" gorm:"primaryKey"`                                             // ID chính (primary key)
	LongURL    string    `json:"long_url" gorm:"column:long_url;type:VARCHAR(512);uniqueIndex"`    // URL gốc (dài), có index unique
	ShortCode  string    `json:"short_code" gorm:"column:short_code;type:varchar(32);uniqueIndex"` // Mã ngắn, có index unique
	CreatedAt  time.Time `json:"created_at"`                                                       // Thời gian tạo
	Clicks     int       `json:"clicks"`                                                           // Số lượt click
}

// Định nghĩa tên bảng trong database
func (URLMapping) TableName() string {
	return "url_mappings"
}
