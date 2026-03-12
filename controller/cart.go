package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/GemaSatya/E-Commerce/model"
)

type CartItemView struct {
	ID       uint
	Name     string
	Price    float64
	Quantity int
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var items []CartItemView
	var total float64

	var cart model.Cart
	if err := model.DB.Where("user_id = ?", currentUser.ID).First(&cart).Error; err == nil {
		var cartItems []model.CartItem
		model.DB.Where("cart_id = ?", cart.ID).Preload("Product").Find(&cartItems)
		for _, item := range cartItems {
			items = append(items, CartItemView{
				ID:       item.ID,
				Name:     item.Product.Name,
				Price:    item.Product.Price,
				Quantity: item.Quantity,
			})
			total += item.Product.Price * float64(item.Quantity)
		}
	}

	t, err := template.ParseFiles("frontend/product/cart.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsLoggedIn bool
		Username   string
		Items      []CartItemView
		Total      float64
	}{
		IsLoggedIn: true,
		Username:   currentUser.Username,
		Items:      items,
		Total:      total,
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "unauthorized"})
		return
	}

	var req struct {
		ProductID uint `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProductID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid request"})
		return
	}

	var product model.Product
	if err := model.DB.First(&product, req.ProductID).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "product not found"})
		return
	}

	// Find or create cart for this user
	var cart model.Cart
	if err := model.DB.Where("user_id = ?", currentUser.ID).First(&cart).Error; err != nil {
		cart = model.Cart{UserID: currentUser.ID}
		if err := model.DB.Create(&cart).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": "failed to create cart"})
			return
		}
	}

	// Upsert cart item: increment quantity if exists, create otherwise
	var cartItem model.CartItem
	if err := model.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		cartItem = model.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  1,
		}
		if err := model.DB.Create(&cartItem).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": "failed to add item"})
			return
		}
	} else {
		cartItem.Quantity++
		if err := model.DB.Save(&cartItem).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": "failed to update item"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "Item added to cart"})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "unauthorized"})
		return
	}

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid item id"})
		return
	}

	var cartItem model.CartItem
	if err := model.DB.First(&cartItem, id).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "item not found"})
		return
	}

	// Verify the cart belongs to the current user before deleting
	var cart model.Cart
	if err := model.DB.Where("id = ? AND user_id = ?", cartItem.CartID, currentUser.ID).First(&cart).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{"error": "forbidden"})
		return
	}

	if err := model.DB.Delete(&cartItem).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "failed to remove item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "Item removed from cart"})
}
