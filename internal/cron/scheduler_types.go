package cron

import (
	"context"
	"sync"
	"time"
)

type Job struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Schedule    string     `json:"schedule"` // Cron expression
	Handler     JobHandler `json:"-"`
	Enabled     bool       `json:"enabled"`
	LastRun     *time.Time `json:"last_run,omitempty"`
	NextRun     *time.Time `json:"next_run,omitempty"`
	RunCount    int        `json:"run_count"`
	Description string     `json:"description"`
}

type JobHandler func(ctx context.Context) error

type Scheduler struct {
	jobs   map[string]*Job
	mu     sync.RWMutex
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}
