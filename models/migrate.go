package models

import (
	"backend/config"
	"log"
)

func Migrate() {

	err := config.DB.AutoMigrate(
		&OTP{},
		&User{},
		&RefreshToken{},
	)

	if err != nil{
		log.Fatal("Migrate is failed !!",err)
	}

	log.Println("All tables  Migrated successfully.")
} 