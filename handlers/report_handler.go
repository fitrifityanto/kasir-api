package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.GetDailyReport(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	dailyReport, err := h.service.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, http.StatusOK, "Successfully retrieved daily report", dailyReport)
}

func (h *ReportHandler) sendResponse(w http.ResponseWriter, status int, message string, data any) {
	response := models.Response{
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
