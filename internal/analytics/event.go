package analytics

// Event represents an analytics event
type Event struct {
	Name        string   `json:"name"`
	Action      string   `json:"action"`
	Label       string   `json:"label"`
	Destinations []string `json:"destinations"`
}

// TrackingDestination represents a destination for analytics events
type TrackingDestination string

const (
	// RedisHLL represents Redis HyperLogLog destination
	RedisHLL TrackingDestination = "redis_hll"
	// Snowplow represents Snowplow destination
	Snowplow TrackingDestination = "snowplow"
)
