package main

import (
	"log"
	"mockapi/internal/application"

	"net/http"

	"github.com/gorilla/mux"
)

type Category struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/categories/{id}", application.GetCategories).Methods("GET")

	log.Println("Server running on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Server failed to start) %v", err)
	}
}
