package models

import (
	"time"
)

// VariableType represents the type of CI variable
type VariableType string

const (
	// VariableTypeEnvVar represents an environment variable
	VariableTypeEnvVar VariableType = "env_var"
	// VariableTypeFile represents a file variable
	VariableTypeFile VariableType = "file"
)

// InstanceVariable represents a CI/CD variable at the instance level
type InstanceVariable struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Key         string      `json:"key" gorm:"uniqueIndex"`
	Value       string      `json:"value"`
	Type        VariableType `json:"variable_type"`
	Description string      `json:"description"`
	Protected   bool        `json:"protected"`
	Masked      bool        `json:"masked"`
	Raw         bool        `json:"raw"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// GetAllInstanceVariables returns all instance variables
func GetAllInstanceVariables() []InstanceVariable {
	var variables []InstanceVariable
	// TODO: Implement database query to get all variables
	return variables
}

// Save saves the instance variable to the database
func (v *InstanceVariable) Save() error {
	// TODO: Implement database save
	return nil
}

// Delete removes the instance variable from the database
func (v *InstanceVariable) Delete() error {
	// TODO: Implement database delete
	return nil
}
