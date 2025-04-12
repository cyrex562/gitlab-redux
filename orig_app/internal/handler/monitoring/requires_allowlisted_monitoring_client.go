package monitoring

import (
	"net"
	"net/http"
)

// RequiresAllowlistedMonitoringClient handles validation of monitoring clients
type RequiresAllowlistedMonitoringClient struct {
	settings Settings
	requestContext RequestContext
	environment Environment
}

// Settings defines the interface for application settings
type Settings interface {
	MonitoringIPWhitelist() []string
	HealthCheckAccessToken() string
}

// RequestContext defines the interface for request context operations
type RequestContext interface {
	ClientIP() string
}

// Environment defines the interface for environment operations
type Environment interface {
	IsDevelopment() bool
}

// NewRequiresAllowlistedMonitoringClient creates a new instance of RequiresAllowlistedMonitoringClient
func NewRequiresAllowlistedMonitoringClient(
	settings Settings,
	requestContext RequestContext,
	environment Environment,
) *RequiresAllowlistedMonitoringClient {
	return &RequiresAllowlistedMonitoringClient{
		settings: settings,
		requestContext: requestContext,
		environment: environment,
	}
}

// ValidateIPAllowlistedOrValidToken validates that the client is either allowlisted by IP or has a valid token
func (r *RequiresAllowlistedMonitoringClient) ValidateIPAllowlistedOrValidToken(w http.ResponseWriter, req *http.Request) bool {
	if r.ClientIPAllowlisted(req) || r.ValidToken(req) {
		return true
	}

	r.Render404(w)
	return false
}

// ClientIPAllowlisted checks if the client IP is allowlisted
func (r *RequiresAllowlistedMonitoringClient) ClientIPAllowlisted(req *http.Request) bool {
	// Always allow developers to access http://localhost:3000/-/metrics for
	// debugging purposes
	if r.environment.IsDevelopment() && req.Host == "localhost:3000" {
		return true
	}

	clientIP := r.requestContext.ClientIP()
	for _, ip := range r.IPAllowlist() {
		if ip.Contains(net.ParseIP(clientIP)) {
			return true
		}
	}

	return false
}

// IPAllowlist returns the list of allowlisted IP addresses
func (r *RequiresAllowlistedMonitoringClient) IPAllowlist() []*net.IPNet {
	allowlist := make([]*net.IPNet, 0)
	for _, ipStr := range r.settings.MonitoringIPWhitelist() {
		_, ipNet, err := net.ParseCIDR(ipStr)
		if err == nil {
			allowlist = append(allowlist, ipNet)
		}
	}
	return allowlist
}

// ValidToken checks if the token is valid
func (r *RequiresAllowlistedMonitoringClient) ValidToken(req *http.Request) bool {
	token := req.URL.Query().Get("token")
	if token == "" {
		token = req.Header.Get("TOKEN")
	}

	if token == "" {
		return false
	}

	return token == r.settings.HealthCheckAccessToken()
}

// Render404 renders a 404 error
func (r *RequiresAllowlistedMonitoringClient) Render404(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
