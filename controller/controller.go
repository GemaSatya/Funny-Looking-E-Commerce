package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/GemaSatya/E-Commerce/model"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
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

	// Parse template dan sisipkan data
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		// fallback ke serve static bila parsing gagal
		http.ServeFile(w, r, "frontend/index.html")
		return
	}

	data := struct{
		ProductsJSON template.JS
	}{
		ProductsJSON: template.JS(b),
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}