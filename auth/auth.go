package auth

import (
	"net/http"

	"github.com/GemaSatya/E-Commerce/model"
)

func RegisterUser(w http.ResponseWriter, r *http.Request){
	// Implementation for user registration
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if !SearchUser(username){
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create new user
	newUser := model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}

	err = model.DB.Create(&newUser).Error
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}