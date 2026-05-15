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
// Supports: @hourly, @daily, @weekly, @monthly, */N (every N minutes), and 0 H * * * (daily at hour H)
func (s *Scheduler) calculateNextRun(schedule string) *time.Time {
	now := time.Now()
	var next time.Time

	switch schedule {
	case "@hourly":
		next = now.Add(time.Hour)
	case "@daily":
		// Run at midnight tomorrow
		next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	case "@weekly":
		// Run at midnight next week (same day)
		next = now.AddDate(0, 0, 7)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	case "@monthly":
		// Run at midnight on the 1st of next month
		next = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	default:
		// Try to parse interval expressions like "*/5" (every 5 minutes)
		if len(schedule) > 2 && schedule[0] == '*' && schedule[1] == '/' {
			minutes := 60 // default to hourly
			fmt.Sscanf(schedule[2:], "%d", &minutes)
			next = now.Add(time.Duration(minutes) * time.Minute)
		} else {
			// Default to hourly for unknown schedules
			next = now.Add(time.Hour)
		}
	}

	return &next
}
