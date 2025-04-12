package commit

import (
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ParseCommitDate handles parsing commit dates
type ParseCommitDate struct {
	logger *service.Logger
}

// NewParseCommitDate creates a new instance of ParseCommitDate
func NewParseCommitDate(
	logger *service.Logger,
) *ParseCommitDate {
	return &ParseCommitDate{
		logger: logger,
	}
}

// ConvertDateToEpoch converts a date string to epoch timestamp
func (p *ParseCommitDate) ConvertDateToEpoch(date string) (int64, error) {
	// If date is empty, return 0
	if date == "" {
		return 0, nil
	}

	// Parse date string
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		p.logger.Error("Failed to parse date", err)
		return 0, err
	}

	// Convert to epoch timestamp
	return t.Unix(), nil
}
