package utils

import (
	"backend/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

//generate Access Token

func GenerateAccessToken(userID uint, Role string) (string, error) {
	Secret_key:=os.Getenv("JWT_ACCESS")
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
		return "","",nil
	}

	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))

	return token,hex.EncodeToString(hash[:]),nil
}

// Save refresh token in db 

func SavaRefreshToken(db *gorm.DB,userId uint,HashedToken string,expiredAt time.Time)error{
	saveReToken:= models.RefreshToken{
		UserId: userId,
		Token: HashedToken,
		ExpiredAt: expiredAt,
	}
	return db.Create(&saveReToken).Error
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

	return &retoken,nil
}

// Delete RefreshToken from DB

func DeleteReToken(db *gorm.DB,token string)error{
	hash:=sha256.Sum256([]byte(token))
	HashedToken := hex.EncodeToString(hash[:])
	return  db.Where("token = ?",HashedToken).Delete(&models.RefreshToken{}).Error
}