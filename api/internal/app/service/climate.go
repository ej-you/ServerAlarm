package service

import (
	"context"
	"log/slog"
	"time"

	"server-alarm/api/internal/app/usecase"
)

// ClimateService represents a service that checks
// temperature with "checkDBEvery" interval.
type ClimateService struct {
	climateUC    *usecase.ClimateUC
	checkDBEvery time.Duration

	started chan struct{}
	exited  chan struct{}
}

// NewClimateService returns a new instance of ClimateService.
func NewClimateService(climateUC *usecase.ClimateUC, checkDBEvery time.Duration) *ClimateService {
	return &ClimateService{
		climateUC:    climateUC,
		checkDBEvery: checkDBEvery,

		started: make(chan struct{}),
		exited:  make(chan struct{}),
	}
}

// Started returns chan that signals the service's readiness.
func (s *ClimateService) Started() <-chan struct{} {
	return s.started
}

// Exited returns chan that signals that service's exited.
func (s *ClimateService) Exited() <-chan struct{} {
	return s.exited
}

// Run starts service. Waits for context done to exit.
func (s *ClimateService) Run(ctx context.Context, errChan chan<- error) {
	slog.Info("start climate service...")
	ticker := time.NewTicker(s.checkDBEvery)
	defer ticker.Stop()

	close(s.started)
	defer close(s.exited)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.climateUC.CheckTemperature(); err != nil {
				errChan <- err
			}
		}
	}
}
