package graphql

import (
	"context"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Result map[string]interface{} `json:"result"`
	More   bool                   `json:"more"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message string `json:"message"`
}

// GraphQLContext represents the context for GraphQL operations
type GraphQLContext struct {
	Channel         interface{}
	CurrentUser     interface{}
	IsSessionlessUser bool
	ScopeValidator   interface{}
}

// GraphQLSchema defines the interface for GraphQL schema operations
type GraphQLSchema interface {
	Execute(ctx context.Context, query string, variables map[string]interface{}, operationName string) (*GraphQLResponse, error)
	DeleteSubscription(subscriptionID string) error
}

// GraphQLChannel represents the WebSocket channel for GraphQL operations
type GraphQLChannel struct {
	subscriptionIDs []string
	schema         GraphQLSchema
	context        *GraphQLContext
}
