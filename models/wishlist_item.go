package models

import "time"

type WishlistItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_product" json:"-"`
	ProductID uint      `gorm:"not null;uniqueIndex:idx_user_product" json:"-"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
