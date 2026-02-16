package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"size:255;not null" json:"name" binding:"required "`
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	Description   string         `gorm:"type:text" json:"description"`
	Category      string         `gorm:"size:255" json:"category"`
	Price         float64        `gorm:"type:decimal(10,2)" json:"price" binding:"required"`
	StockQuantity uint           `gorm:"not null;default:0" json:"stock_quantity" binding:"required"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
