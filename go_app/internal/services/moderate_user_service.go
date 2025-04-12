package services

import (
	"errors"

	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
)

// ModerateUserService handles user moderation based on abuse reports
type ModerateUserService struct {
	report   *models.AbuseReport
	user     *models.User
	params   models.AbuseReportParams
}

// NewModerateUserService creates a new instance of ModerateUserService
func NewModerateUserService(report *models.AbuseReport, user *models.User, params models.AbuseReportParams) *ModerateUserService {
	return &ModerateUserService{
		report: report,
		user:   user,
		params: params,
	}
}

// Execute performs the moderation operation
func (s *ModerateUserService) Execute() (*ServiceResponse, error) {
	if s.report == nil {
		return nil, errors.New("report is required")
	}

	if s.user == nil {
		return nil, errors.New("user is required")
	}

	// Get the reported user
	reportedUser, err := models.GetUserByID(s.report.UserID)
	if err != nil {
		return nil, err
	}

	// Perform the moderation action based on user_action
	switch s.params.UserAction {
	case "ban":
		if err := reportedUser.Ban(s.user); err != nil {
			return nil, err
		}
	case "block":
		if err := reportedUser.Block(s.user); err != nil {
			return nil, err
		}
	case "warn":
		if err := reportedUser.Warn(s.user, s.params.Reason); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid user action")
	}

	// Update the report status
	s.report.Status = models.AbuseReportStatusClosed
	s.report.Reason = s.params.Reason
	s.report.Comment = s.params.Comment

	// Save the changes
	if err := s.report.Save(); err != nil {
		return nil, err
	}

	return &ServiceResponse{
		Success: true,
		Message: "User moderated successfully",
	}, nil
}
