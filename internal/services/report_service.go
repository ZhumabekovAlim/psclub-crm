package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) SummaryReport(ctx context.Context, from, to time.Time) (*models.SummaryReport, error) {
	return s.repo.SummaryReport(ctx, from, to)
}
func (s *ReportService) AdminsReport(ctx context.Context, from, to time.Time) (*models.AdminsReport, error) {
	return s.repo.AdminsReport(ctx, from, to)
}
func (s *ReportService) SalesReport(ctx context.Context, from, to time.Time) (*models.SalesReport, error) {
	return s.repo.SalesReport(ctx, from, to)
}
func (s *ReportService) AnalyticsReport(ctx context.Context, from, to time.Time) (*models.AnalyticsReport, error) {
	return s.repo.AnalyticsReport(ctx, from, to)
}
func (s *ReportService) DiscountsReport(ctx context.Context, from, to time.Time) (*models.DiscountsReport, error) {
	return s.repo.DiscountsReport(ctx, from, to)
}
