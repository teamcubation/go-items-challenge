package utils

import (
	"encoding/json"
	"net/http"
)

func RespondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		http.Error(w, "Failed to encode error message", http.StatusInternalServerError)
	}
}
