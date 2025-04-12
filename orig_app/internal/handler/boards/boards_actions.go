package boards

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// BoardsActions handles board-related actions
type BoardsActions struct {
	boardCreateService *service.BoardCreateService
	boardFinder        *service.BoardFinder
	boardVisitService  *service.BoardVisitService
	authService        *service.AuthService
	logger             *util.Logger
}

// NewBoardsActions creates a new instance of BoardsActions
func NewBoardsActions(
	boardCreateService *service.BoardCreateService,
	boardFinder *service.BoardFinder,
	boardVisitService *service.BoardVisitService,
	authService *service.AuthService,
	logger *util.Logger,
) *BoardsActions {
	return &BoardsActions{
		boardCreateService: boardCreateService,
		boardFinder:        boardFinder,
		boardVisitService:  boardVisitService,
		authService:        authService,
		logger:             logger,
	}
}

// Index handles the index action for boards
func (b *BoardsActions) Index(ctx *gin.Context) error {
	// Authorize read board
	if err := b.authorizeReadBoard(ctx); err != nil {
		return err
	}

	// Redirect to recent board if applicable
	if err := b.redirectToRecentBoard(ctx); err != nil {
		return err
	}

	// Get the board
	board, err := b.getBoard(ctx)
	if err != nil {
		return err
	}

	// If no board exists, create one
	if board == nil {
		result, err := b.boardCreateService.Execute(ctx)
		if err != nil {
			return err
		}
		board = result.Payload
	}

	// Push licensed features
	b.pushLicensedFeatures(ctx)

	// Render the index view
	ctx.HTML(200, "boards/index.html", gin.H{
		"board": board,
	})

	return nil
}

// Show handles the show action for a board
func (b *BoardsActions) Show(ctx *gin.Context) error {
	// Authorize read board
	if err := b.authorizeReadBoard(ctx); err != nil {
		return err
	}

	// Get the board
	board, err := b.getBoard(ctx)
	if err != nil {
		return err
	}

	if board == nil {
		return util.NewNotFoundError("board not found")
	}

	// Add / update the board in the recent visits table
	parent := b.getParent(ctx)
	currentUser := b.getCurrentUser(ctx)
	if err := b.boardVisitService.New(parent, currentUser).Execute(ctx, board); err != nil {
		return err
	}

	// Push licensed features
	b.pushLicensedFeatures(ctx)

	// Render the show view
	ctx.HTML(200, "boards/show.html", gin.H{
		"board": board,
	})

	return nil
}

// authorizeReadBoard authorizes read access to the board
func (b *BoardsActions) authorizeReadBoard(ctx *gin.Context) error {
	parent := b.getParent(ctx)
	return b.authService.AuthorizeReadBoard(ctx, parent)
}

// redirectToRecentBoard redirects to the most recently visited board if applicable
func (b *BoardsActions) redirectToRecentBoard(ctx *gin.Context) error {
	parent := b.getParent(ctx)
	if !parent.MultipleIssueBoardsAvailable() {
		return nil
	}

	latestVisitedBoard, err := b.getLatestVisitedBoard(ctx)
	if err != nil {
		return err
	}

	if latestVisitedBoard == nil {
		return nil
	}

	// Redirect to the latest visited board
	boardPath := b.getBoardPath(ctx, latestVisitedBoard.Board)
	ctx.Redirect(302, boardPath)
	return nil
}

// getLatestVisitedBoard gets the latest visited board
func (b *BoardsActions) getLatestVisitedBoard(ctx *gin.Context) (*model.BoardVisit, error) {
	parent := b.getParent(ctx)
	currentUser := b.getCurrentUser(ctx)
	return b.boardVisitService.NewFinder(parent, currentUser).Latest(ctx)
}

// pushLicensedFeatures pushes licensed features to the frontend
// Noop on FOSS
func (b *BoardsActions) pushLicensedFeatures(ctx *gin.Context) {
	// This is a noop in the FOSS version
}

// getBoard gets the board
func (b *BoardsActions) getBoard(ctx *gin.Context) (*model.Board, error) {
	boards, err := b.boardFinder.Execute(ctx)
	if err != nil {
		return nil, err
	}

	if len(boards) == 0 {
		return nil, nil
	}

	return boards[0], nil
}

// getParent gets the parent (group or project)
func (b *BoardsActions) getParent(ctx *gin.Context) model.Parent {
	if b.isGroup(ctx) {
		return b.getGroup(ctx)
	}
	return b.getProject(ctx)
}

// getBoardPath gets the path for a board
func (b *BoardsActions) getBoardPath(ctx *gin.Context, board *model.Board) string {
	parent := b.getParent(ctx)
	if b.isGroup(ctx) {
		return b.getGroupBoardPath(parent, board)
	}
	return b.getProjectBoardPath(parent, board)
}

// isGroup checks if the parent is a group
func (b *BoardsActions) isGroup(ctx *gin.Context) bool {
	_, exists := ctx.Get("group")
	return exists
}

// Helper methods for getting context values
func (b *BoardsActions) getCurrentUser(ctx *gin.Context) *model.User {
	user, _ := ctx.Get("current_user")
	return user.(*model.User)
}

func (b *BoardsActions) getGroup(ctx *gin.Context) *model.Group {
	group, _ := ctx.Get("group")
	return group.(*model.Group)
}

func (b *BoardsActions) getProject(ctx *gin.Context) *model.Project {
	project, _ := ctx.Get("project")
	return project.(*model.Project)
}

// Path helper methods
func (b *BoardsActions) getGroupBoardPath(group model.Parent, board *model.Board) string {
	return "/groups/" + group.GetID() + "/boards/" + board.ID
}

func (b *BoardsActions) getProjectBoardPath(project model.Parent, board *model.Board) string {
	return "/projects/" + project.GetID() + "/boards/" + board.ID
}
