package controller

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func LoadTemplate(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "frontend/index.html")
}