package application

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Category struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := []Category{
		{Name: "sports", Active: true},
		{Name: "electronics", Active: true},
		{Name: "books", Active: true},
		{Name: "fashion", Active: true},
		{Name: "toys", Active: false},
		{Name: "furniture", Active: true},
		{Name: "music", Active: true},
		{Name: "movies", Active: false},
		{Name: "games", Active: true},
		{Name: "outdoors", Active: false},
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid category ID"}`, http.StatusBadRequest)
		return
	}

	if id < 0 || id >= len(categories) {
		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories[id]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
