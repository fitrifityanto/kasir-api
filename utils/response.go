package utils

import (
	"encoding/json"
	"kasir-api/models"
	"net/http"
)

func MessageResponse(w http.ResponseWriter, status int, message string, data any) {
	response := models.Response{
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
