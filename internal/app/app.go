package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rps/config"
	"rps/internal/db"
	"rps/internal/db/store"
	router "rps/internal/http"
	"rps/internal/http/handler"
	"rps/internal/logging"
	"rps/internal/metric"
	"rps/internal/observability"
	"rps/internal/service"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config *config.Config
	logger *slog.Logger
	server *http.Server
	db     *pgxpool.Pool
	obs    *observability.Observability
}

func newHttpServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func New() (a *App, err error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("errors load config: %w", err)
	}

	// logger
	logger := logging.New(cfg.Log)

	a = &App{
		config: cfg,
		logger: logger,
	}

	defer func() {
		if err != nil && a != nil {
			if shutdownErr := a.shutdown(); shutdownErr != nil {
				logger.Error(
					"startup cleanup failed",
					"error",
					shutdownErr,
				)
			}
		}
	}()

	// db
	logger.Info("connecting database")
	a.db, err = db.NewPostgres(cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	// observability
	logger.Info("initializing observability")
	initCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.obs, err = observability.New(initCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("setup observability: %w", err)
	}

	logger.Info("building router")
	// metrics
	metrics := metric.NewMetrics()

	// handler
	handlers, err := buildDependencies(a.db, logger, metrics)
	if err != nil {
		return nil, fmt.Errorf("build handlers: %w", err)
	}

	// route
	r := router.NewRouter(logger, handlers, metrics)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	a.server = newHttpServer(addr, r)

	return a, nil

}

func buildDependencies(db *pgxpool.Pool, logger *slog.Logger, m *metric.Metrics) (*handler.Handlers, error) {
	st := store.New(db)
	policyService := service.NewPolicyService(st, logger)
	return &handler.Handlers{
		Policy: handler.NewPolicyHandler(policyService),
		Metric: m.NewHttpHandler(),
	}, nil
}

func (a *App) Start() error {
	errCh := make(chan error, 1)

	go a.startHTTPServer(errCh)

	return a.wait(errCh)
}

func (a *App) startHTTPServer(errCh chan<- error) {
	a.logger.Info("starting HTTP server", "addr", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		errCh <- fmt.Errorf("http server: %w", err)
	}
}

func (a *App) wait(errCh <-chan error) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-errCh:
		a.logger.Error("runtime error, shutting down", "error", err)
		shutdownErr := a.shutdown()
		return errors.Join(err, shutdownErr)

	case <-ctx.Done():
		a.logger.Info("Shutdown signal received")
		return a.shutdown()
	}
}
func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errs []error

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown http server: %w", err))
		}
	}

	if a.db != nil {
		a.db.Close()
	}

	if a.obs != nil {
		if err := a.obs.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown observability: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	a.logger.Info("application stopped gracefully")
	return nil
}
