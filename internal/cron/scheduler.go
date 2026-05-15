package cron

import (
	"context"
	"fmt"
	"log"
	"time"
)

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:   make(map[string]*Job),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Minute) // Check every minute
	go s.run()
}

func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.cancel()
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			s.checkJobs()
		}
	}
}

func (s *Scheduler) checkJobs() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()

	for _, job := range s.jobs {
		if !job.Enabled {
			continue
		}

		// Check if job should run now
		if job.NextRun != nil && now.After(*job.NextRun) {
			go s.executeJob(job)
		}
	}
}

func (s *Scheduler) executeJob(job *Job) {
	log.Printf("Executing job: %s", job.Name)

	err := job.Handler(s.ctx)
	now := time.Now()

	s.mu.Lock()
	job.LastRun = &now
	job.RunCount++

	// Calculate next run time based on schedule
	job.NextRun = s.calculateNextRun(job.Schedule)
	s.mu.Unlock()

	if err != nil {
		log.Printf("Job %s failed: %v", job.Name, err)
	} else {
		log.Printf("Job %s completed successfully", job.Name)
	}
}

// calculateNextRun calculates the next run time based on a cron expression
// For simplicity, this is a basic implementation that supports hourly/daily schedules
func (s *Scheduler) calculateNextRun(schedule string) *time.Time {
	now := time.Now()
	var next time.Time

	switch schedule {
	case "@hourly":
		next = now.Add(time.Hour)
	case "@daily":
		next = now.Add(24 * time.Hour)
	case "@weekly":
		next = now.Add(7 * 24 * time.Hour)
	default:
		// Default to hourly for unknown schedules
		next = now.Add(time.Hour)
	}

	return &next
}

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job.ID == "" {
		job.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Calculate initial next run time
	job.NextRun = s.calculateNextRun(job.Schedule)

	s.jobs[job.ID] = job
	return nil
}

// RemoveJob removes a job from the scheduler
func (s *Scheduler) RemoveJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[jobID]; !exists {
		return fmt.Errorf("job not found")
	}

	delete(s.jobs, jobID)
	return nil
}

// GetJob returns a specific job
func (s *Scheduler) GetJob(jobID string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

// GetJobs returns all jobs
func (s *Scheduler) GetJobs() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// EnableJob enables a job
func (s *Scheduler) EnableJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found")
	}

	job.Enabled = true
	job.NextRun = s.calculateNextRun(job.Schedule)
	return nil
}

// DisableJob disables a job
func (s *Scheduler) DisableJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found")
	}

	job.Enabled = false
	job.NextRun = nil
	return nil
}

// RunJobNow manually runs a job immediately
func (s *Scheduler) RunJobNow(jobID string) error {
	s.mu.RLock()
	job, exists := s.jobs[jobID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job not found")
	}

	go s.executeJob(job)
	return nil
}

// DefaultJobs returns default scheduled jobs for infrastructure management
func DefaultJobs(cleanupManager any) []*Job {
	// These would be implemented with actual cleanup manager integration
	return []*Job{
		{
			Name:        "dev-environment-cleanup",
			Schedule:    "@hourly",
			Description: "Check and cleanup dev environments",
			Enabled:     true,
			Handler: func(ctx context.Context) error {
				// This would call the cleanup manager to check and cleanup dev environments
				log.Println("Running dev environment cleanup check")
				return nil
			},
		},
		{
			Name:        "test-environment-cleanup",
			Schedule:    "@hourly",
			Description: "Check and cleanup test environments",
			Enabled:     true,
			Handler: func(ctx context.Context) error {
				// This would call the cleanup manager to check and cleanup test environments
				log.Println("Running test environment cleanup check")
				return nil
			},
		},
		{
			Name:        "cost-alert-check",
			Schedule:    "@daily",
			Description: "Check cost alerts and notify",
			Enabled:     true,
			Handler: func(ctx context.Context) error {
				// This would check cost thresholds and send alerts
				log.Println("Running cost alert check")
				return nil
			},
		},
	}
}
