package auth

import (
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