package analytics

import (
	"context"
)

// AnalyticsService exposes metrics to the frontend
type AnalyticsService struct {
	ctx    context.Context
	engine *AnalyticsEngine
}

func NewAnalyticsService(engine *AnalyticsEngine) *AnalyticsService {
	return &AnalyticsService{
		engine: engine,
	}
}

func (a *AnalyticsService) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *AnalyticsService) GetStats() AnalyticsSnapshot {
	return a.engine.GetSnapshot()
}
