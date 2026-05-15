package cron

import (
	"context"
	"fmt"
	"log"
	"time"
)

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job.ID == "" {
		job.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

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
	return []*Job{
		{
			Name:        "dev-environment-cleanup",
			Schedule:    "@hourly",
			Description: "Check and cleanup dev environments",
			Enabled:     true,
			Handler: func(ctx context.Context) error {
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
				log.Println("Running cost alert check")
				return nil
			},
		},
	}
}
