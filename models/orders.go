package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	UserID        uint    `gorm:"not null;index" json:"user_id"`
	TotalAmount   float64 `gorm:"not null" json:"total_amount"`
	Address       string  `gorm:"type:text;not null" json:"address"`
	Status        string  `gorm:"type:varchar(20);default:'pending'" json:"status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User       User        `gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}
