package controller

import (
	"encoding/json"
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
	json.NewEncoder(w).Encode(map[string]interface{}{
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)

}