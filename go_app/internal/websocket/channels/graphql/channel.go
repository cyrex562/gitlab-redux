package graphql

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/gitlab-org/gitlab/internal/websocket/channel"
	"gitlab.com/gitlab-org/gitlab/internal/websocket/logging"
)

// NewGraphQLChannel creates a new GraphQLChannel instance
func NewGraphQLChannel(channel *channel.Channel, schema GraphQLSchema) *GraphQLChannel {
	return &GraphQLChannel{
		schema: schema,
		context: &GraphQLContext{
			Channel:         channel,
			IsSessionlessUser: false,
		},
	}
}

// Subscribe handles the subscription request
func (c *GraphQLChannel) Subscribe(ctx context.Context, params map[string]interface{}) error {
	// Parse the GraphQL request
	request := &GraphQLRequest{}
	if err := parseRequest(params, request); err != nil {
		return fmt.Errorf("invalid request: %v", err)
	}

	// Execute the GraphQL query
	result, err := c.schema.Execute(ctx, request.Query, request.Variables, request.OperationName)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	// Track subscription ID if present
	if subscriptionID, ok := result.Result["subscription_id"].(string); ok {
		c.subscriptionIDs = append(c.subscriptionIDs, subscriptionID)
	}

	// Send the response
	response := &GraphQLResponse{
		Result: result.Result,
		More:   result.More,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	// Send the response through the channel
	if err := c.context.Channel.(*channel.Channel).Send(responseJSON); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	// Log the subscription
	logging.GetLogger().Info(logging.LogPayload{
		Params: map[string]interface{}{
			"operation_name": request.OperationName,
			"has_subscription": len(c.subscriptionIDs) > 0,
		},
	}, "GraphQL subscription established")

	return nil
}

// Unsubscribe handles the unsubscription request
func (c *GraphQLChannel) Unsubscribe(ctx context.Context) error {
	if len(c.subscriptionIDs) == 0 {
		return nil
	}

	// Delete all subscriptions
	for _, subscriptionID := range c.subscriptionIDs {
		if err := c.schema.DeleteSubscription(subscriptionID); err != nil {
			logging.GetLogger().Error(logging.LogPayload{
				Params: map[string]interface{}{
					"subscription_id": subscriptionID,
				},
			}, "Failed to delete subscription", err)
		}
	}

	// Log the unsubscription
	logging.GetLogger().Info(logging.LogPayload{
		Params: map[string]interface{}{
			"subscription_count": len(c.subscriptionIDs),
		},
	}, "GraphQL subscriptions cleaned up")

	c.subscriptionIDs = nil
	return nil
}

// parseRequest parses the request parameters into a GraphQLRequest
func parseRequest(params map[string]interface{}, request *GraphQLRequest) error {
	// Parse query
	if query, ok := params["query"].(string); ok {
		request.Query = query
	} else {
		return fmt.Errorf("query is required")
	}

	// Parse variables
	if variables, ok := params["variables"].(map[string]interface{}); ok {
		request.Variables = variables
	} else {
		request.Variables = make(map[string]interface{})
	}

	// Parse operation name
	if operationName, ok := params["operationName"].(string); ok {
		request.OperationName = operationName
	}

	return nil
}
