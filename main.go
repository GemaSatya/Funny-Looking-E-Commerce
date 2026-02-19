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

	// Product Routes
	mux.HandleFunc("/add-product", controller.AddProduct)	// Post Product
	mux.HandleFunc("/product/{id}", controller.GetOneProduct)	// Get One Product
	mux.HandleFunc("/product/delete/{id}", controller.DeleteProduct) // Delete Product
	
	fmt.Println("Server is running on port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Server failed to start:", err)
	}

}