package services

import (
	"backend/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

//generate 6 DIGIT OTP 

func GenerateOTP() (string,error) {
	n, err := rand.Int(rand.Reader,big.NewInt(900000))
	if err != nil{
		return "",err
	}
	otp := 100000 + n.Int64()
	return fmt.Sprintf("%06d",otp),nil
}

//OTP hashing

func HashOTP(otp string)string{
	hash := sha256.Sum256([]byte(otp))
	return  hex.EncodeToString(hash[:])
}

//create OTP with hashing save DB

func CreateOTP(db *gorm.DB,userID uint,purpose string,expiryMinutes int) (string,error){
	otp ,err := GenerateOTP()
	if err != nil{
		return  "",err
	}

	otpHash := HashOTP(otp)

	tx :=db.Begin()

	if err :=tx.Model(&models.OTP{}).Where("user_id = ? AND purpose = ? AND is_used = false ",userID,purpose).
	Update("is_used",true).Error;err !=nil{
		tx.Rollback()
		return  "",err
	}

	updatedVersion  := models.OTP{
		UserID: userID,
		OTPCode: otpHash,
		Purpose: purpose,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(expiryMinutes)),
		IsUsed: false,
	}

	if err := tx.Create(&updatedVersion).Error; err != nil {
		tx.Rollback()
		return  "",err
	}

	tx.Commit()
	return otp,nil
}


//verify OTP

func VerifyOTP(db *gorm.DB,userID uint,purpose string,otp string) error{
	hasedOTP := HashOTP(otp)

	var record models.OTP

	if err := db.Where(`user_id = ? AND purpose = ? AND otp_code = ? AND is_used = false AND expires_at > ?`,userID,purpose,hasedOTP,time.Now()).
	First(&record).Error;
	 err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return errors.New("Invalid or expired OTP")
		} 
		return  err
	}

	return db.Model(&record).Update("is_used",true).Error
}
