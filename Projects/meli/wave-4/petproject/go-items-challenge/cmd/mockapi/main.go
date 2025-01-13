package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teamcubation/go-items-challenge/internal/application"
)

type Category struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

//func getCategories(w http.ResponseWriter, r *http.Request) {
//	// Lista hardcodeada de categorías
//	categories := []Category{
//		{Name: "sports", Active: true},
//		{Name: "electronics", Active: true},
//		{Name: "books", Active: true},
//		{Name: "fashion", Active: true},
//		{Name: "toys", Active: false},
//		{Name: "furniture", Active: true},
//		{Name: "music", Active: true},
//		{Name: "movies", Active: false},
//		{Name: "games", Active: true},
//		{Name: "outdoors", Active: false},
//	}
//
//	vars := mux.Vars(r)
//	idStr := vars["id"]
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		http.Error(w, `{"error": "Invalid category ID"}`, http.StatusBadRequest)
//		return
//	}
//
//	// Validar el rango del ID
//	if id < 0 || id >= len(categories) {
//		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(categories[id])
//}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/categories/{id}", application.GetCategories).Methods("GET")

	//srv := &http.Server{
	//	Addr:    ":8000",
	//	Handler: r,
	//}

	log.Println("Server running on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Server failed to start) %v", err)
	}
}
