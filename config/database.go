package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil{
		panic("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s password=%s user=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_USER"),)

	database,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil{
		panic("Database connection is failed ")
	}
	DB=database

	fmt.Print("Database connection connect successfully")

}

