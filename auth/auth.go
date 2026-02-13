package auth

import (
	"net/http"
	"time"

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

func LoginUser(w http.ResponseWriter, r *http.Request){
	// Implementation for user login
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	if SearchUser(username){
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	var user model.User
	var login model.Login

	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(password, user.Password){
		http.Error(w, "Wrong Credentials", http.StatusUnauthorized)
		return
	}

	if !SearchToken(user.ID){
		http.Error(w, "User already logged in", http.StatusConflict)
		return
	}

	sessionToken := GenerateToken(32)
	csrfToken := GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name : "session_token",
		Value : sessionToken,
		Expires : time.Now().Add(30 * time.Second),
		HttpOnly : true,
	})

	http.SetCookie(w, &http.Cookie{
		Name : "csrf_token",
		Value : csrfToken,
		Expires : time.Now().Add(30 * time.Second),
		HttpOnly : false,
	})

	login = model.Login{
		HashedPassword: user.Password,
		SessionToken:   sessionToken,
		CSRFToken:      csrfToken,
		SessionId:      user.ID,
	}

	err := model.DB.Create(&login).Error
	if err != nil {
		http.Error(w, "Error creating login session", http.StatusInternalServerError)
		return
	}
}