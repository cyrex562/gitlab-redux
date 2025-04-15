package services

import (
	"errors"

	"github.com/cyrex562/gitlab-redux/internal/models"
)

// ServiceResponse represents a response from a service
type ServiceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateAbuseReportService handles updating abuse reports
type UpdateAbuseReportService struct {
	report   *models.AbuseReport
	user     *models.User
	params   models.AbuseReportParams
}

// NewUpdateAbuseReportService creates a new instance of UpdateAbuseReportService
func NewUpdateAbuseReportService(report *models.AbuseReport, user *models.User, params models.AbuseReportParams) *UpdateAbuseReportService {
	return &UpdateAbuseReportService{
		report: report,
		user:   user,
		params: params,
	}
}

// Execute performs the update operation
func (s *UpdateAbuseReportService) Execute() (*ServiceResponse, error) {
	if s.report == nil {
		return nil, errors.New("report is required")
	}

	if s.user == nil {
		return nil, errors.New("user is required")
	}

	// Update report fields
	if s.params.Close {
		s.report.Status = models.AbuseReportStatusClosed
	}
	s.report.Reason = s.params.Reason
	s.report.Comment = s.params.Comment
	s.report.LabelIDs = s.params.LabelIDs

	// Save the changes
	if err := s.report.Save(); err != nil {
		return nil, err
	}

	return &ServiceResponse{
		Success: true,
		Message: "Abuse report updated successfully",
	}, nil
}
