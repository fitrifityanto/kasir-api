package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"log"
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
		h.sendResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	dailyReport, err := h.service.GetDailyReport()
	if err != nil {
		log.Printf("Failed to get daily report: %v", err)

		h.sendResponse(w, http.StatusInternalServerError, "Internal sever error", nil)
		return
	}

	message := "Successfully retrieved daily report"
	if dailyReport.TotalSales == 0 {
		message = "No transaction recorded today"

	}

	h.sendResponse(w, http.StatusOK, message, dailyReport)
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
