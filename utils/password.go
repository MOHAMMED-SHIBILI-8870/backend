package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return  string(hash),err
}

func ComparePassword(password string,hash string)bool{
	err :=  bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return  err == nil
}