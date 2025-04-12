package services

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// SidekiqService handles interaction with Sidekiq for background job monitoring
type SidekiqService struct {
	redisClient *redis.Client
}

// NewSidekiqService creates a new instance of SidekiqService
func NewSidekiqService(redisClient *redis.Client) *SidekiqService {
	return &SidekiqService{
		redisClient: redisClient,
	}
}

// Stats represents Sidekiq statistics
type Stats struct {
	Processed    int64     `json:"processed"`
	Failed       int64     `json:"failed"`
	Scheduled    int64     `json:"scheduled"`
	Retries      int64     `json:"retries"`
	Dead         int64     `json:"dead"`
	Processes    int64     `json:"processes"`
	Workers      int64     `json:"workers"`
	Enqueued     int64     `json:"enqueued"`
	Busy         int64     `json:"busy"`
	LastUpdateAt time.Time `json:"last_update_at"`
}

// Queue represents a Sidekiq queue
type Queue struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Latency  int64  `json:"latency"`
	Paused   bool   `json:"paused"`
}

// Worker represents a Sidekiq worker
type Worker struct {
	ID        string    `json:"id"`
	Hostname  string    `json:"hostname"`
	StartedAt time.Time `json:"started_at"`
	Queues    []string  `json:"queues"`
	Busy      int64     `json:"busy"`
}

// Job represents a Sidekiq job
type Job struct {
	ID        string                 `json:"id"`
	Class     string                 `json:"class"`
	Args      []interface{}          `json:"args"`
	Queue     string                 `json:"queue"`
	CreatedAt time.Time              `json:"created_at"`
	EnqueuedAt time.Time             `json:"enqueued_at"`
	Error     string                 `json:"error,omitempty"`
	Retry     bool                   `json:"retry"`
	RetryCount int64                 `json:"retry_count,omitempty"`
	RetriedAt  time.Time             `json:"retried_at,omitempty"`
	FailedAt   time.Time             `json:"failed_at,omitempty"`
	Backtrace  []string              `json:"backtrace,omitempty"`
}

// GetStats retrieves Sidekiq statistics
func (s *SidekiqService) GetStats() (*Stats, error) {
	ctx := context.Background()

	// Get processed count
	processed, err := s.redisClient.Get(ctx, "stat:processed").Int64()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// Get failed count
	failed, err := s.redisClient.Get(ctx, "stat:failed").Int64()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// Get scheduled count
	scheduled, err := s.redisClient.ZCard(ctx, "schedule").Int64()
	if err != nil {
		return nil, err
	}

	// Get retries count
	retries, err := s.redisClient.ZCard(ctx, "retry").Int64()
	if err != nil {
		return nil, err
	}

	// Get dead count
	dead, err := s.redisClient.ZCard(ctx, "dead").Int64()
	if err != nil {
		return nil, err
	}

	// Get processes count
	processes, err := s.redisClient.SCard(ctx, "processes").Int64()
	if err != nil {
		return nil, err
	}

	// Get workers count
	workers, err := s.redisClient.SCard(ctx, "workers").Int64()
	if err != nil {
		return nil, err
	}

	// Get enqueued count
	enqueued, err := s.redisClient.LLen(ctx, "queue:default").Int64()
	if err != nil {
		return nil, err
	}

	// Get busy count
	busy, err := s.redisClient.SCard(ctx, "busy").Int64()
	if err != nil {
		return nil, err
	}

	return &Stats{
		Processed:    processed,
		Failed:       failed,
		Scheduled:    scheduled,
		Retries:      retries,
		Dead:         dead,
		Processes:    processes,
		Workers:      workers,
		Enqueued:     enqueued,
		Busy:         busy,
		LastUpdateAt: time.Now(),
	}, nil
}

// GetQueues retrieves information about all Sidekiq queues
func (s *SidekiqService) GetQueues() ([]*Queue, error) {
	ctx := context.Background()

	// Get all queue names
	queueNames, err := s.redisClient.SMembers(ctx, "queues").Result()
	if err != nil {
		return nil, err
	}

	queues := make([]*Queue, 0, len(queueNames))
	for _, name := range queueNames {
		// Get queue size
		size, err := s.redisClient.LLen(ctx, "queue:"+name).Int64()
		if err != nil {
			return nil, err
		}

		// Get queue latency
		latency, err := s.redisClient.Get(ctx, "queue:"+name+":latency").Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		// Check if queue is paused
		paused, err := s.redisClient.SIsMember(ctx, "paused_queues", name).Result()
		if err != nil {
			return nil, err
		}

		queues = append(queues, &Queue{
			Name:    name,
			Size:    size,
			Latency: latency,
			Paused:  paused,
		})
	}

	return queues, nil
}

// GetWorkers retrieves information about all Sidekiq workers
func (s *SidekiqService) GetWorkers() ([]*Worker, error) {
	ctx := context.Background()

	// Get all worker IDs
	workerIDs, err := s.redisClient.SMembers(ctx, "workers").Result()
	if err != nil {
		return nil, err
	}

	workers := make([]*Worker, 0, len(workerIDs))
	for _, id := range workerIDs {
		// Get worker info
		info, err := s.redisClient.HGetAll(ctx, "worker:"+id).Result()
		if err != nil {
			return nil, err
		}

		startedAt, err := time.Parse(time.RFC3339, info["started_at"])
		if err != nil {
			return nil, err
		}

		// Get worker queues
		queues, err := s.redisClient.SMembers(ctx, "worker:"+id+":queues").Result()
		if err != nil {
			return nil, err
		}

		// Get worker busy count
		busy, err := s.redisClient.SCard(ctx, "worker:"+id+":busy").Int64()
		if err != nil {
			return nil, err
		}

		workers = append(workers, &Worker{
			ID:        id,
			Hostname:  info["hostname"],
			StartedAt: startedAt,
			Queues:    queues,
			Busy:      busy,
		})
	}

	return workers, nil
}

// GetScheduledJobs retrieves all scheduled jobs
func (s *SidekiqService) GetScheduledJobs() ([]*Job, error) {
	return s.getJobsFromSortedSet("schedule")
}

// GetRetries retrieves all retry jobs
func (s *SidekiqService) GetRetries() ([]*Job, error) {
	return s.getJobsFromSortedSet("retry")
}

// GetDeadJobs retrieves all dead jobs
func (s *SidekiqService) GetDeadJobs() ([]*Job, error) {
	return s.getJobsFromSortedSet("dead")
}

// Helper function to get jobs from a sorted set
func (s *SidekiqService) getJobsFromSortedSet(setName string) ([]*Job, error) {
	ctx := context.Background()

	// Get all job IDs from the set
	jobIDs, err := s.redisClient.ZRange(ctx, setName, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0, len(jobIDs))
	for _, id := range jobIDs {
		// Get job data
		data, err := s.redisClient.HGetAll(ctx, setName+":"+id).Result()
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339, data["created_at"])
		if err != nil {
			return nil, err
		}

		enqueuedAt, err := time.Parse(time.RFC3339, data["enqueued_at"])
		if err != nil {
			return nil, err
		}

		job := &Job{
			ID:         id,
			Class:      data["class"],
			Queue:      data["queue"],
			CreatedAt:  createdAt,
			EnqueuedAt: enqueuedAt,
			Error:      data["error"],
			Retry:      data["retry"] == "true",
		}

		if retryCount, err := strconv.ParseInt(data["retry_count"], 10, 64); err == nil {
			job.RetryCount = retryCount
		}

		if retriedAt, err := time.Parse(time.RFC3339, data["retried_at"]); err == nil {
			job.RetriedAt = retriedAt
		}

		if failedAt, err := time.Parse(time.RFC3339, data["failed_at"]); err == nil {
			job.FailedAt = failedAt
		}

		if backtrace, err := s.redisClient.LRange(ctx, setName+":"+id+":backtrace", 0, -1).Result(); err == nil {
			job.Backtrace = backtrace
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
