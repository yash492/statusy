package applications

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/port/httphandler"
	"github.com/yash492/statusy/internal/port/httphandler/generated/api"
)

type ServerDeps struct {
	lg      *slog.Logger
	readDB  *pgxpool.Pool
	writeDB *pgxpool.Pool
}

func NewServerDeps(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) ServerDeps {
	return ServerDeps{
		lg:      lg,
		readDB:  readDB,
		writeDB: writeDB,
	}
}

type ServerApplication struct {
	HttpHandler httphandler.Handler
	lg          *slog.Logger
}

func NewServerApplication(deps ServerDeps) ServerApplication {
	lg := deps.lg
	servicesRepo := servicesdb.NewPostgresServiceRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	incidentsRepo := incidentsdb.NewPostgresIncidentRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	handler := httphandler.Handler{
		ListStatuspageCmd:       command.NewListStatuspageCmd(lg, servicesRepo),
		IncidentByStatuspageCmd: command.NewIncidentByStatuspageCmd(lg, servicesRepo, incidentsRepo),
	}
	return ServerApplication{
		HttpHandler: handler,
		lg:          lg,
	}
}

func (s ServerApplication) Start(ctx context.Context, addr string) error {
	serverInterface := api.NewStrictHandler(s.HttpHandler, nil)
	r := chi.NewRouter()
	handler := api.HandlerFromMux(serverInterface, r)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	s.lg.Info("starting http server", slog.String("addr", addr))

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.lg.Error("http server stopped", slog.Any("err", err))
			return err
		}
		return nil
	case <-ctx.Done():
		s.lg.Info("shutdown signal received, stopping http server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.lg.Error("graceful shutdown failed", slog.Any("err", err))
			return err
		}

		s.lg.Info("http server stopped gracefully")
		return nil
	}
}
