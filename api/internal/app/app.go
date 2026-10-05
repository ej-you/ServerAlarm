package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"server-alarm/api/config"
	"server-alarm/api/internal/app/repo"
	"server-alarm/api/internal/app/service"
	"server-alarm/api/internal/app/usecase"
	"server-alarm/api/internal/pkg/db"
	"server-alarm/api/internal/pkg/healthcheck"
	"server-alarm/api/internal/pkg/logger"
	"server-alarm/api/internal/pkg/ntfy"
	"server-alarm/api/internal/pkg/storage"
)

// Service describes app service started in parallel with other similar services.
type Service interface {
	// Started returns chan that signals that service's readiness.
	Started() <-chan struct{}
	// Exited returns chan that signals that service's exited.
	Exited() <-chan struct{}
	// Run starts service. Waits to context done to exit.
	Run(ctx context.Context, errChan chan<- error)
}

// App represents full application with a list of service.
type App struct {
	services []Service
}

// New returns a new instance of App.
func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}
	// setup logger
	logOpts := []logger.SlogOption{logger.WithLogLevel(cfg.Logging.LogLevel)}
	if cfg.Logging.JSONFormat {
		logOpts = append(logOpts, logger.WithJSONFormat())
	}
	logger.InitSlog(logOpts...)
	// print out config params in debug log mode
	slog.Debug("current config", "config", cfg)
	slog.Info("checking config",
		"check db every", cfg.App.CheckDBEvery.String(),
		"temperature treshold", cfg.App.TempTreshold)

	// init sources
	storageInst := storage.NewKeyValueInMem()
	dbInst, err := db.New(cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("create db conn: %w", err)
	}
	ntfyClient, err := ntfy.NewClient(ntfy.WithTokenAuth(cfg.Ntfy.Token),
		ntfy.WithHTTP(),
		ntfy.WithHost(cfg.Ntfy.Host),
		ntfy.WithPort(cfg.Ntfy.Port))
	if err != nil {
		return nil, fmt.Errorf("create ntfy client: %w", err)
	}
	// init repos
	climateRepoDB := repo.NewClimateRepoDB(dbInst)
	climateRepoCache := repo.NewClimateRepoCache(storageInst)
	climateRepoNtfy := repo.NewClimateRepoNtfy(ntfyClient, cfg.App.Locale, cfg.Ntfy.Theme)
	// init usecases
	climateUC := usecase.NewClimateUC(
		climateRepoCache,
		climateRepoDB,
		climateRepoNtfy,
		cfg.App.Climate.TempTreshold.Standart,
	)

	healthcheckService := healthcheck.New([]healthcheck.Checking{dbInst},
		healthcheck.WithPort(cfg.App.HealthCheck.Port))
	climateService := service.NewClimateService(climateUC, cfg.App.CheckDBEvery)

	slog.Info("init app", "status", "OK")
	return &App{
		[]Service{healthcheckService, climateService},
	}, nil
}

// Run starts app. This function is blocking.
// It waits for os signal to gracefully shutdown.
func (a *App) Run() (appErr error) {
	// ctx for app services
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle shutdown process signals
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// wait groups to sync starting/exiting all services
	var wgStarted, wgExited sync.WaitGroup

	// start services
	serviceErr := make(chan error, 1)
	for _, service := range a.services {
		wgStarted.Go(func() {
			<-service.Started()
		})
		wgExited.Go(func() {
			<-service.Exited()
		})
		go service.Run(appContext, serviceErr)
	}

	wgStarted.Wait()
	select {
	case handledSignal := <-quitSig:
		cancel()
		slog.Warn("got os signal; shutdown app...", "signal", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		slog.Warn("one of services fell down; shutdown app...")
	case <-appContext.Done():
		appErr = appContext.Err()
		slog.Warn("context canceled; shutdown app...")
	}
	wgExited.Wait()
	slog.Info("all services exited")
	return appErr
}
