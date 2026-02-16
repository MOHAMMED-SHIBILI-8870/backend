package models

import (
	"time"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID  uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint          `gorm:"not null;index" json:"product_id"`
	Product   Product       `gorm:"foreignKey:ProductID"`
	UnitPrice float64        `gorm:"not null" json:"unit_price"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	TotalPrice float64       `gorm:"not null" json:"total_price"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
