package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GemaSatya/E-Commerce/auth"
	"github.com/GemaSatya/E-Commerce/controller"
	"github.com/GemaSatya/E-Commerce/model"
	"github.com/joho/godotenv"
)

func main(){

	mux := http.NewServeMux()

	err := godotenv.Load()
	if err != nil{
		fmt.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")

	model.ConnectDatabase()

	// mux.HandleFunc("/", controller.HelloHandler)
	mux.HandleFunc("/", controller.LoadTemplate)

	// Auth Routes
	mux.HandleFunc("/register", auth.RegisterUser)	// Register User
	mux.HandleFunc("/login", auth.LoginUser)			// Login User
	mux.HandleFunc("/logout", auth.LogoutUser)		// Logout User

	// Product Routes
	mux.HandleFunc("/add-product", controller.AddProduct)			// Post Product
	mux.HandleFunc("/product/{id}", controller.GetOneProduct)		// Get One Product
	mux.HandleFunc("/product/delete/{id}", controller.DeleteProduct)	// Delete Product

	// Cart Routes
	mux.HandleFunc("/cart", controller.ViewCart)				// View Cart
	mux.HandleFunc("/cart/add", controller.AddToCart)			// Add Item to Cart
	mux.HandleFunc("/cart/item/{id}", controller.RemoveFromCart)		// Remove Item from Cart

	fmt.Println("Server is running on port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Server failed to start:", err)
	}

}