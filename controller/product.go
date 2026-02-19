package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/GemaSatya/E-Commerce/model"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Implementation for adding a product
	var productRequest struct{
		Name string
		Description string
		Price float64
	}

	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product := model.Product{
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
	}

	if err := model.DB.Create(&product).Error; err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Product added successfully",
		"product": product,
	})

}

func GetOneProduct(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil{
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product model.Product
	if err := model.DB.First(&product, id).Error; err != nil{
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Parse template dan serve HTML dengan data produk
	t, err := template.ParseFiles("frontend/product/product.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, product); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

}

func DeleteProduct(w http.ResponseWriter, r *http.Request){
	// Implementation for deleting a product
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil{
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product model.Product
	if err := model.DB.First(&product, id).Error; err != nil{
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err := model.DB.Delete(&product).Error; err != nil{
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Product deleted successfully",
	})

}