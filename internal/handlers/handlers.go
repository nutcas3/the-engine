package handlers

import (
	"context"
	"time"

	"the-engine/internal/alerts"
	"the-engine/internal/cache"
	"the-engine/internal/cleanup"
	"the-engine/internal/cron"
	"the-engine/internal/health"
	"the-engine/internal/kubernetes"
	"the-engine/internal/nuke"
	"the-engine/internal/rate"
)

// Handlers holds the dependencies for HTTP handlers
type Handlers struct {
	k8sClient      *kubernetes.Client
	healthChecker  *health.Checker
	cache          *cache.Cache
	RateLimiter    *rate.RateLimiter
	alertManager   *alerts.AlertManager
	cleanupManager *cleanup.CleanupManager
	cronScheduler  *cron.Scheduler
	nukeManager    *nuke.NukeManager
}

// NewHandlers creates a new Handlers instance
func NewHandlers(k8sClient *kubernetes.Client) *Handlers {
	h := &Handlers{
		k8sClient:      k8sClient,
		healthChecker:  health.NewChecker("1.0.0"),
		cache:          cache.NewCache(5 * time.Minute),
		RateLimiter:    rate.NewRateLimiter(100, 10),
		alertManager:   alerts.NewAlertManager(),
		cleanupManager: cleanup.NewCleanupManager(),
		cronScheduler:  cron.NewScheduler(),
		nukeManager:    nuke.NewNukeManager(),
	}

	h.alertManager.Start(context.Background())
	h.cleanupManager.Start(context.Background())
	h.cronScheduler.Start()

	return h
}
