package members

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// MembersPresentation handles the presentation of members
type MembersPresentation struct {
	membersPreloader *service.MembersPreloader
	presenterFactory *service.PresenterFactory
	logger           *service.Logger
}

// NewMembersPresentation creates a new instance of MembersPresentation
func NewMembersPresentation(
	membersPreloader *service.MembersPreloader,
	presenterFactory *service.PresenterFactory,
	logger *service.Logger,
) *MembersPresentation {
	return &MembersPresentation{
		membersPreloader: membersPreloader,
		presenterFactory: presenterFactory,
		logger:           logger,
	}
}

// PresentMembers presents a list of members
func (m *MembersPresentation) PresentMembers(ctx *gin.Context, members []*service.Member) ([]*service.MemberPresenter, error) {
	// Preload associations for the members
	if err := m.preloadAssociations(members); err != nil {
		return nil, err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil, err
	}

	// Create and fabricate presenters
	presenters, err := m.presenterFactory.Fabricate(
		members,
		map[string]interface{}{
			"current_user": currentUser,
			"presenter_class": "MembersPresenter",
		},
	)
	if err != nil {
		return nil, err
	}

	// Convert presenters to MemberPresenter type
	memberPresenters := make([]*service.MemberPresenter, len(presenters))
	for i, presenter := range presenters {
		if memberPresenter, ok := presenter.(*service.MemberPresenter); ok {
			memberPresenters[i] = memberPresenter
		}
	}

	return memberPresenters, nil
}

// preloadAssociations preloads all associations for the members
func (m *MembersPresentation) preloadAssociations(members []*service.Member) error {
	return m.membersPreloader.PreloadAll(members)
}
