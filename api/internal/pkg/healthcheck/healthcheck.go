package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Checking describes item for health check.
type Checking interface {
	// IsReady checks that item is ready to use. Returns nil error if item is ready.
	IsReady() error
}

// HealthCheck represents a light HTTP-server with a
// single route to check all items health.
type HealthCheck struct {
	// default: 0.0.0.0
	host string
	// default: 80
	port string
	// default: /health
	route string

	items []Checking

	started chan struct{}
	exited  chan struct{}
}

type Option func(*HealthCheck)

// New returns a new instance of HealthCheck.
func New(items []Checking, options ...Option) *HealthCheck {
	inst := &HealthCheck{
		host:    "0.0.0.0",
		port:    "80",
		route:   "/health",
		items:   items,
		started: make(chan struct{}),
		exited:  make(chan struct{}),
	}
	for _, option := range options {
		option(inst)
	}
	return inst
}

// WithHost sets custom host for HealthCheck server.
func WithHost(host string) Option {
	return func(h *HealthCheck) {
		h.host = host
	}
}

// WithPort sets custom port for HealthCheck server.
func WithPort(port string) Option {
	return func(h *HealthCheck) {
		h.port = port
	}
}

// WithRoute sets custom route for HealthCheck server.
func WithRoute(route string) Option {
	return func(h *HealthCheck) {
		h.route = route
	}
}

// Started returns chan that signals the service's readiness.
func (h *HealthCheck) Started() <-chan struct{} {
	return h.started
}

// Exited returns chan that signals that service's exited.
func (h *HealthCheck) Exited() <-chan struct{} {
	return h.exited
}

// Run starts server. Waits to context done to exit.
func (h *HealthCheck) Run(ctx context.Context, errChan chan<- error) {
	addr := net.JoinHostPort(h.host, h.port)
	addrFull := "http://" + addr
	slog.Info("start healthcheck service...", "addr", addrFull, "route", h.route)

	// listen port
	listener, err := net.Listen("tcp", addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- fmt.Errorf("start healthcheck listener: %w", err)
		return
	}

	// init server
	mux := http.NewServeMux()
	mux.HandleFunc(h.route, h.healthHanlder)
	srv := &http.Server{Addr: addr, Handler: mux}

	// serve connections
	go func() {
		err := srv.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("start healthcheck server: %w", err)
		}
	}()
	close(h.started)

	// gracefully shutdown server
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Warn("shutdown healthcheck server", "err", err)
	}
	close(h.exited)
}

// healthHanlder handles HealthCheck route.
func (h *HealthCheck) healthHanlder(w http.ResponseWriter, r *http.Request) {
	var (
		health = true
		err    error
	)
	// check items
	for _, item := range h.items {
		if itemErr := item.IsReady(); itemErr != nil {
			health = false
			err = multierror.Append(err, itemErr)
		}
	}
	// return response
	if !health {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
