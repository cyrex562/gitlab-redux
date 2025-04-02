package projects

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/go/internal/activitypub"
	"gitlab.com/gitlab-org/gitlab/go/internal/activitypub/projects/services"
	"gitlab.com/gitlab-org/gitlab/go/internal/activitypub/serializers"
	"gitlab.com/gitlab-org/gitlab/go/internal/models"
)

// ReleasesController handles ActivityPub endpoints for project releases
type ReleasesController struct {
	ApplicationController
}

// Index handles GET requests to list the releases actor
func (c *ReleasesController) Index(w http.ResponseWriter, r *http.Request) error {
	project := c.GetProject()
	opts := map[string]string{
		"inbox":  c.inboxURL(project),
		"outbox": c.outboxURL(project),
	}

	serializer := serializers.NewReleasesActorSerializer()
	response := serializer.Represent(project, opts)

	return json.NewEncoder(w).Encode(response)
}

// Inbox handles POST requests to the releases inbox
func (c *ReleasesController) Inbox(w http.ResponseWriter, r *http.Request) error {
	if err := c.enforcePayload(w, r); err != nil {
		return err
	}

	service := c.inboxService()
	success := true
	var err error

	if service != nil {
		success, err = service.Execute()
		if err != nil {
			return err
		}
	}

	response := map[string]interface{}{
		"success": success,
	}

	if !success && service != nil {
		response["errors"] = service.Errors()
	}

	return json.NewEncoder(w).Encode(response)
}

// Outbox handles GET requests to list the releases outbox
func (c *ReleasesController) Outbox(w http.ResponseWriter, r *http.Request) error {
	releases := c.releases(r.URL.Query())
	serializer := serializers.NewReleasesOutboxSerializer().WithPagination(r, w)
	response := serializer.Represent(releases)

	return json.NewEncoder(w).Encode(response)
}

// releases returns a list of releases for the project
func (c *ReleasesController) releases(params map[string][]string) []*models.Release {
	finder := models.NewReleasesFinder(c.GetProject(), c.GetCurrentUser(), params)
	return finder.Execute()
}

// enforcePayload ensures the request has a valid JSON payload
func (c *ReleasesController) enforcePayload(w http.ResponseWriter, r *http.Request) error {
	if _, err := c.payload(r); err == nil {
		return nil
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	return nil
}

// payload parses and returns the request body as JSON
func (c *ReleasesController) payload(r *http.Request) (map[string]interface{}, error) {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// isFollow checks if the payload is a Follow activity
func (c *ReleasesController) isFollow(payload map[string]interface{}) bool {
	return payload["type"] == "Follow"
}

// isUnfollow checks if the payload is an Undo activity of a Follow
func (c *ReleasesController) isUnfollow(payload map[string]interface{}) bool {
	if payload["type"] != "Undo" {
		return false
	}

	object, ok := payload["object"].(map[string]interface{})
	if !ok {
		return false
	}

	return object["type"] == "Follow"
}

// inboxService returns the appropriate service based on the payload type
func (c *ReleasesController) inboxService() services.Service {
	payload, err := c.payload(r)
	if err != nil {
		return nil
	}

	if c.isFollow(payload) {
		return services.NewReleasesFollowService(c.GetProject(), payload)
	}

	if c.isUnfollow(payload) {
		return services.NewReleasesUnfollowService(c.GetProject(), payload)
	}

	return nil
}

// inboxURL returns the URL for the project releases inbox
func (c *ReleasesController) inboxURL(project *models.Project) string {
	return activitypub.URLFor("inbox_project_releases", project)
}

// outboxURL returns the URL for the project releases outbox
func (c *ReleasesController) outboxURL(project *models.Project) string {
	return activitypub.URLFor("outbox_project_releases", project)
}
