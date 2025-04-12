package service

import (
	"context"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// OnboardingService handles onboarding-related operations
type OnboardingService struct {
	db *model.DB
}

// NewOnboardingService creates a new instance of OnboardingService
func NewOnboardingService(db *model.DB) *OnboardingService {
	return &OnboardingService{
		db: db,
	}
}

// GetOnboardingStatusPresenter gets the onboarding status presenter
func (s *OnboardingService) GetOnboardingStatusPresenter(ctx context.Context) *OnboardingStatusPresenter {
	// Get the current user
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return &OnboardingStatusPresenter{
			user: nil,
		}
	}

	// Get the user's invites
	invites, err := s.getUserInvites(ctx, user)
	if err != nil {
		return &OnboardingStatusPresenter{
			user:   user,
			invites: []*model.Member{},
		}
	}

	return &OnboardingStatusPresenter{
		user:   user,
		invites: invites,
	}
}

// getCurrentUser gets the current user
func (s *OnboardingService) getCurrentUser(ctx context.Context) (*model.User, error) {
	// TODO: Implement current user retrieval
	// This should:
	// 1. Get the current user from the context
	// 2. Return the user
	return nil, nil
}

// getUserInvites gets the user's invites
func (s *OnboardingService) getUserInvites(ctx context.Context, user *model.User) ([]*model.Member, error) {
	// TODO: Implement user invites retrieval
	// This should:
	// 1. Get the user's invites from the database
	// 2. Return the invites
	return []*model.Member{}, nil
}

// OnboardingStatusPresenter presents the onboarding status
type OnboardingStatusPresenter struct {
	user    *model.User
	invites []*model.Member
}

// IsSingleInvite checks if there's a single invite
func (p *OnboardingStatusPresenter) IsSingleInvite() bool {
	return len(p.invites) == 1
}

// GetLastInvitedMember gets the last invited member
func (p *OnboardingStatusPresenter) GetLastInvitedMember() *model.Member {
	if len(p.invites) == 0 {
		return nil
	}
	return p.invites[len(p.invites)-1]
}

// GetLastInvitedMemberSource gets the last invited member source
func (p *OnboardingStatusPresenter) GetLastInvitedMemberSource() interface{} {
	member := p.GetLastInvitedMember()
	if member == nil {
		return nil
	}
	return member.Source
}
