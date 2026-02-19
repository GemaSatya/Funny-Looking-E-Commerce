package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/GemaSatya/E-Commerce/model"
)

func getCurrentUser(r *http.Request) *model.User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	var login model.Login
	if err := model.DB.Where("session_token = ?", cookie.Value).First(&login).Error; err != nil {
		return nil
	}

	var user model.User
	if err := model.DB.First(&user, login.SessionId).Error; err != nil {
		return nil
	}

	return &user
}

func LoadTemplate(w http.ResponseWriter, r *http.Request) {
	// Ambil semua produk dari database
	var products []model.Product
	if err := model.DB.Find(&products).Error; err != nil {
		// bila error, gunakan list kosong supaya halaman masih bisa dirender
		products = []model.Product{}
	}

	// Marshal ke JSON untuk disuntikkan ke frontend
	b, err := json.Marshal(products)
	if err != nil {
		b = []byte("[]")
	}

	// Get current user
	currentUser := getCurrentUser(r)
	var username string
	if currentUser != nil {
		username = currentUser.Username
	}

	// Parse template dan sisipkan data
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		// fallback ke serve static bila parsing gagal
		http.ServeFile(w, r, "frontend/index.html")
		return
	}

	data := struct{
		ProductsJSON template.JS
		Username     string
		IsLoggedIn   bool
	}{
		ProductsJSON: template.JS(b),
		Username:     username,
		IsLoggedIn:   currentUser != nil,
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}