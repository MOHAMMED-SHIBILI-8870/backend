package seeder

import (
	"backend/config"
	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

func AdminSeeder(db *gorm.DB)error{
	var user models.User

	err := config.DB.Where("email = ?","admin@gmail.com").First(&user).Error

	if err == nil{
		return nil
	}

	hashPass,err := utils.HashPassword("admin@123")

	if err != nil {
		return err
	}

	admin := models.User{
		FullName: "admin",
		Email: "admin@gmail.com",
		HashPassword: hashPass,
		Role: "admin",
	}

	return  db.Create(&admin).Error
}