package main

import (
	"fmt"
	"net/http"
	"os"

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

	mux.HandleFunc("/", controller.HelloHandler)
	mux.HandleFunc("/frontend", controller.LoadTemplate)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Server failed to start:", err)
	}

}