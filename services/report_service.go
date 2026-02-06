package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDailyReport() (*models.DailyReport, error) {
	return s.repo.GetDailyReport()
}

func (s *ReportService) GetReport(startDate, endDate string) (*models.FullReport, error) {
	return s.repo.GetReport(startDate, endDate)
}
