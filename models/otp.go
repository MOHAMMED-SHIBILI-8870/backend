package models

import "time"

type OTP struct {
	ID        uint      `gorm:"primaryKey;autoIncreament" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	OTPCode   string    `gorm:"type:varchar(255);not null" json:"otp_code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Purpose   string    `gorm:"type:varchar(50);not null" json:"purpose"`
	IsUsed    bool      `gorm:"default:false;not null" json:"is_used"`
	CreatedAt time.Time 
}
