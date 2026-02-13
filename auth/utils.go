package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/GemaSatya/E-Commerce/model"
	"golang.org/x/crypto/bcrypt"
)

func SearchUser(username string) bool {
	var user model.User
	
	err := model.DB.Where("username = ?", username).First(&user).Error

	return err != nil

}

func HashPassword(password string) (string, error) {
	// Dummy hash function for illustration
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SearchToken(sessionId uint) bool{
	var user model.Login

	err := model.DB.Where("session_id = ?", sessionId).First(&user).Error

	return err != nil
}

func GenerateToken(length int) string{
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil{
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}