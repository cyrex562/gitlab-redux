package services

import (
	"errors"

	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
)

// UpdateInstanceVariablesService handles updating CI variables
type UpdateInstanceVariablesService struct {
	variables []models.InstanceVariable
	user      *models.User
}

// NewUpdateInstanceVariablesService creates a new instance of UpdateInstanceVariablesService
func NewUpdateInstanceVariablesService(variables []models.InstanceVariable, user *models.User) *UpdateInstanceVariablesService {
	return &UpdateInstanceVariablesService{
		variables: variables,
		user:      user,
	}
}

// Execute performs the update operation
func (s *UpdateInstanceVariablesService) Execute() error {
	if s.user == nil {
		return errors.New("user is required")
	}

	// Get existing variables
	existingVariables := models.GetAllInstanceVariables()
	existingMap := make(map[uint]*models.InstanceVariable)
	for i := range existingVariables {
		existingMap[existingVariables[i].ID] = &existingVariables[i]
	}

	// Process updates and deletions
	for _, variable := range s.variables {
		if variable.ID == 0 {
			// New variable
			if err := variable.Save(); err != nil {
				return err
			}
		} else {
			// Update existing variable
			if existing, ok := existingMap[variable.ID]; ok {
				existing.Key = variable.Key
				existing.Value = variable.Value
				existing.Type = variable.Type
				existing.Description = variable.Description
				existing.Protected = variable.Protected
				existing.Masked = variable.Masked
				existing.Raw = variable.Raw
				if err := existing.Save(); err != nil {
					return err
				}
				delete(existingMap, variable.ID)
			}
		}
	}

	// Delete removed variables
	for _, variable := range existingMap {
		if err := variable.Delete(); err != nil {
			return err
		}
	}

	return nil
}
