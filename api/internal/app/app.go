package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server-alarm/api/config"
	"server-alarm/api/internal/app/repo"
	"server-alarm/api/internal/app/usecase"
	"server-alarm/api/internal/pkg/db"
	"server-alarm/api/internal/pkg/logger"
	"server-alarm/api/internal/pkg/ntfy"
	"server-alarm/api/internal/pkg/storage"
)

// App represents full application with a single service.
type App struct {
	climateUC    *usecase.ClimateUC
	checkDBEvery time.Duration
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
		cfg.App.TempTreshold,
	)

	slog.Info("init app", "status", "OK")

	return &App{
		climateUC:    climateUC,
		checkDBEvery: cfg.App.CheckDBEvery,
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

	// start service
	serviceErr := make(chan error, 1)
	go a.runService(appContext, serviceErr)

	select {
	case handledSignal := <-quitSig:
		cancel()
		slog.Warn("got os signal; shutdown app...", "signal", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		slog.Warn("service fell down; shutdown app...")
	case <-appContext.Done():
		appErr = appContext.Err()
		slog.Warn("context canceled; shutdown app...")
	}
	return appErr
}

// runService starts app service.
// Service checks temperature with "checkDBEvery" interval.
func (a *App) runService(ctx context.Context, errChan chan<- error) {
	slog.Info("start check temperature service", "status", "OK")
	ticker := time.NewTicker(a.checkDBEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			//
			if err := a.climateUC.CheckTemperature(); err != nil {
				errChan <- err
			}
		}
	}
}
