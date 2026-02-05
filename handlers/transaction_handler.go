package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strings"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item apa aja, quantitty nya
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		h.sendResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if len(req.Items) == 0 {

		h.sendResponse(w, http.StatusBadRequest, "Items cannot be empty", nil)
		return
	}

	products, err := h.service.Checkout(req.Items)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "stock") {
			h.sendResponse(w, http.StatusBadRequest, err.Error(), nil)
		} else {
			h.sendResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		}
		return
	}

	h.sendResponse(w, http.StatusCreated, "Successfully checkout", products)
}

func (h *TransactionHandler) sendResponse(w http.ResponseWriter, status int, message string, data any) {
	response := models.Response{
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
