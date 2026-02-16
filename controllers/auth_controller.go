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
		Fullname string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
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
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

func VerifyOTPController(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		OTP     string `json:"otp" binding:"required,len=6"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	valid, err := services.VerifyOTP(user.ID, input.OTP, input.Purpose)

	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP is expired or wrong",
		})
		return
	}

	if input.Purpose == "signup" {
		config.DB.Model(&user).Updates(map[string]interface{}{
			"is_verified": true,
			"updated_at":  time.Now(),
		})
	}

	config.DB.Where("user_id = ? AND purpose = ?", user.ID, input.Purpose).
		Delete(&models.OTP{})

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
	})

}

func Login(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var users models.User

	err := config.DB.Where("email = ?", input.Email).First(&users).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "this Email doesn't match",
		})
		return
	}

	if !users.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user doesn't verified with OTP.",
		})
		return
	}

	if users.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "user was blocked by Admin",
		})
		return
	}

	if !utils.ComparePassword(input.Password, users.HashPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "password is not correct",
		})
		return
	}

	accessToken, err := utils.GenerateAccessToken(users.ID, users.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "could not generating access token",
		})
		return
	}

	refreshToken, hashedToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generating  refresh token",
		})
		return
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	err = utils.SaveRefreshToken(config.DB, users.ID, hashedToken, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", accessToken, 20*120, "/", "", false, false)
	c.SetCookie("refresh_token", refreshToken, int(time.Until(expiresAt).Seconds()), "/", "", false, true)

	c.JSON(200, gin.H{
		"status":       "your Logged in",
		"role":         users.Role,
		"access_token": accessToken,
	})
}

func ForgetPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, "user not found")
		return
	}

	otp, err := services.CreateOTP(config.DB, user.ID, "reset_password", 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate OTP",
		})
		return
	}

	if err := services.SentOTPEmail(input.Email, otp, "reset_password"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not send OTP email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "OTP sent to your email for password reset",
	})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Email      string `json:"email" binding:"required,email"`
		NewPasword string `json:"new_password" binding:"required,min=4"`
		OTP        string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	valid, err := services.VerifyOTP(user.ID, input.OTP, "reset_password")

	if !valid || err != nil {
		c.String(http.StatusBadRequest, "Invalid or expired token")
		return
	}

	hashedPass, err := utils.HashPassword(input.NewPasword)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := config.DB.Model(models.User{}).Where("email = ?", user.Email).
		Updates(map[string]interface{}{"hash_password": hashedPass, "updated_at": time.Now()}).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	config.DB.Where("user_id = ? AND purpose=?", user.ID, "reset_password").Delete(&models.OTP{})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Password reset successfully",
	})
}

func ResendOtpHandler(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	if user.IsVerified && input.Purpose == "signup" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already verified",
		})
		return
	}

	otp, err := services.CreateOTP(config.DB, user.ID, input.Purpose, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate OTP",
		})
		return
	}

	if err := services.SentOTPEmail(input.Email, otp, input.Purpose); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not send OTP email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP resent successfully. Please check your email",
	})
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	if err := utils.DeleteReToken(config.DB, refreshToken); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)


	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
