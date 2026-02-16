package utils

import (
	"backend/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

//generate Access Token

func GenerateAccessToken(userID uint, Role string) (string, error) {
	Secret_key:=os.Getenv("JWT_SECRETKEY")
	claims := jwt.MapClaims{
		"user_id":userID,
		"role":Role,
		"exp":time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return  token.SignedString([]byte(Secret_key))
}

//create random token plain and hashed

func GenerateRefreshToken()(string,string,error){
	b:=make([]byte,32)
	_,err := rand.Read(b)
	if err != nil{
		return "","",err
	}

	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))

	return token,hex.EncodeToString(hash[:]),nil
}

// Save refresh token in db 

func SaveRefreshToken(db *gorm.DB, userId uint, hashedToken string,expiresAt time.Time) error {

	db.Where("user_id = ?", userId).Delete(&models.RefreshToken{})

	refreshToken := models.RefreshToken{
		UserId:    userId,
		Token:     hashedToken,
		ExpiredAt: expiresAt,
	}

	fmt.Println(refreshToken)

	return db.Create(&refreshToken).Error
}


//Validate refresh token

func ValidateRefreshToken(db *gorm.DB,token string)(*models.RefreshToken,error){
	hash := sha256.Sum256([]byte(token))
	hashtoken:=hex.EncodeToString(hash[:])


	var retoken models.RefreshToken

	err := db.Where("token=? AND expires_at > ?",hashtoken,time.Now()).First(&retoken).Error
	if err != nil {
		return nil,errors.New("expired or invalid Refresh token")
	}

	return &retoken,err
}

// Delete RefreshToken from DB

func DeleteReToken(db *gorm.DB,token string)error{
	hash:=sha256.Sum256([]byte(token))
	HashedToken := hex.EncodeToString(hash[:])
	return  db.Where("token = ?",HashedToken).Delete(&models.RefreshToken{}).Error
}