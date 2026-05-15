package cron

import (
	"context"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	s := NewScheduler()
	if s == nil {
		t.Fatal("Expected non-nil scheduler")
	}
	if s.jobs == nil {
		t.Error("Expected jobs map to be initialized")
	}
	if s.ctx == nil {
		t.Error("Expected context to be initialized")
	}
}

func TestScheduler_AddJob(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	err := s.AddJob(job)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if job.ID == "" {
		t.Error("Expected job ID to be set")
	}
	
	if job.NextRun == nil {
		t.Error("Expected NextRun to be set")
	}
}

func TestScheduler_RemoveJob(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job)
	
	err := s.RemoveJob(job.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Remove non-existent job
	err = s.RemoveJob("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent job")
	}
}

func TestScheduler_GetJob(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job)
	
	retrieved, err := s.GetJob(job.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.Name != "test-job" {
		t.Errorf("Expected test-job, got %s", retrieved.Name)
	}
	
	// Get non-existent job
	_, err = s.GetJob("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent job")
	}
}

func TestScheduler_GetJobs(t *testing.T) {
	s := NewScheduler()
	
	job1 := &Job{
		Name:     "job1",
		Schedule: "@hourly",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	job2 := &Job{
		Name:     "job2",
		Schedule: "@daily",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job1)
	s.AddJob(job2)
	
	jobs := s.GetJobs()
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}
}

func TestScheduler_EnableJob(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Enabled:  false,
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job)
	
	err := s.EnableJob(job.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	retrieved, _ := s.GetJob(job.ID)
	if !retrieved.Enabled {
		t.Error("Expected job to be enabled")
	}
}

func TestScheduler_DisableJob(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Enabled:  true,
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job)
	
	err := s.DisableJob(job.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	retrieved, _ := s.GetJob(job.ID)
	if retrieved.Enabled {
		t.Error("Expected job to be disabled")
	}
	if retrieved.NextRun != nil {
		t.Error("Expected NextRun to be nil for disabled job")
	}
}

func TestScheduler_RunJobNow(t *testing.T) {
	s := NewScheduler()
	
	job := &Job{
		Name:     "test-job",
		Schedule: "@hourly",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	
	s.AddJob(job)
	
	err := s.RunJobNow(job.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Run non-existent job
	err = s.RunJobNow("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent job")
	}
}

func TestScheduler_calculateNextRun(t *testing.T) {
	s := NewScheduler()
	
	tests := []struct {
		schedule string
	}{
		{"@hourly"},
		{"@daily"},
		{"@weekly"},
		{"@monthly"},
		{"*/5"},
		{"*/15"},
	}
	
	for _, tt := range tests {
		nextRun := s.calculateNextRun(tt.schedule)
		if nextRun == nil {
			t.Errorf("Expected non-nil NextRun for schedule %s", tt.schedule)
		}
		if nextRun.Before(time.Now()) {
			t.Errorf("Expected NextRun to be in the future for schedule %s", tt.schedule)
		}
	}
}

func TestDefaultJobs(t *testing.T) {
	jobs := DefaultJobs(nil)
	if len(jobs) != 3 {
		t.Errorf("Expected 3 default jobs, got %d", len(jobs))
	}
	
	for _, job := range jobs {
		if job.Name == "" {
			t.Error("Expected job name to be set")
		}
		if job.Schedule == "" {
			t.Error("Expected job schedule to be set")
		}
		if !job.Enabled {
			t.Error("Expected default jobs to be enabled")
		}
	}
}
