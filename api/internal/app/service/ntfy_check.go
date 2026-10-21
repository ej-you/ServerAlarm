package service

import (
	"context"
	"log/slog"
	"time"

	"server-alarm/api/internal/app/usecase"
)

// NtfyCheckService represents a service that checks
// ntfy service availability by sending message.
type NtfyCheckService struct {
	ntfyCheckUC *usecase.NtfyCheckUC
	period      time.Duration

	started chan struct{}
	exited  chan struct{}
}

// NewNtfyCheckService returns a new instance of NtfyCheckService.
func NewNtfyCheckService(ntfyCheckUC *usecase.NtfyCheckUC,
	period time.Duration) *NtfyCheckService {

	return &NtfyCheckService{
		ntfyCheckUC: ntfyCheckUC,
		period:      period,

		started: make(chan struct{}),
		exited:  make(chan struct{}),
	}
}

// Started returns chan that signals the service's readiness.
func (s *NtfyCheckService) Started() <-chan struct{} {
	return s.started
}

// Exited returns chan that signals that service's exited.
func (s *NtfyCheckService) Exited() <-chan struct{} {
	return s.exited
}

// Run starts service. Waits for context done to exit.
func (s *NtfyCheckService) Run(ctx context.Context, errChan chan<- error) {
	slog.Info("start ntfy check service...")
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	close(s.started)
	defer close(s.exited)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.ntfyCheckUC.SendCheckMsg(); err != nil {
				errChan <- err
			}
		}
	}
}
