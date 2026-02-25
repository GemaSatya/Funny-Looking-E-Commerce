package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/GemaSatya/E-Commerce/model"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Serve HTML form for GET requests
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("frontend/product/add-product.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle POST requests for adding products
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var productRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	// Parse JSON from request body
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Invalid JSON format",
		})
		return
	}

	// Validate input
	if productRequest.Name == "" || productRequest.Description == "" || productRequest.Price <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Invalid input: name, description, and price (>0) are required",
		})
		return
	}

	product := model.Product{
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
	}

	if err := model.DB.Create(&product).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Failed to create product",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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