package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserId    uint   `gorm:"not null"`
	Token     string `gorm:"not null;unique"`
	ExpiredAt time.Time `gorm:"not null"`
	CreateAt time.Time 
	UpdatedAt time.Time
	DeleteAt time.Time `gorm:"index"`
}