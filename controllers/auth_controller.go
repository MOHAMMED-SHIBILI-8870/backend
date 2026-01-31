package controllers

import (
	"backend/config"
	"backend/models"
	"backend/services"
	"backend/utils"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		Fullname string `json:"full_name" gorm:"not null"`
		Email    string `json:"email" gorm:"not null"`
		Password string `json:"password" gorm:"not null"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existUser models.User

	if err := config.DB.Where("email = ?", input.Email).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	hashpass, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating account",
		})
		return
	}

	user := models.User{
		FullName:     input.Fullname,
		Email:        input.Email,
		HashPassword: hashpass,
		Role:         "user",
		IsVerifed:    false,
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create user",
		})
		return
	}

	otp, err := services.CreateOTP(config.DB, user.ID, "signup", 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate OTP",
		})
		return
	}

	if err := services.SentOTPEmail(input.Email, otp, "signup"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not send OTP email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered  successfully.OTP sent to your email.",
	})
}

func VerifyOTP(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		OTP     string `json:"otp" binding:"required,len=6"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?",input.Email).First(&user).Error;err != nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error":"user not found",
		})
		return
	}


	if err := services.VerifyOTP(config.DB,user.ID,input.Purpose,input.OTP);err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"failed to verify user", 
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"OTP verified successfully",
	})

}
