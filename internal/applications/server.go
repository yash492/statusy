package applications

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/scheduledmaintenancesdb"
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
	scheduledMaintenancesRepo := scheduledmaintenancesdb.NewPostgresScheduledMaintenanceRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	handler := httphandler.Handler{
		ListStatuspageCmd:                   command.NewListStatuspageCmd(lg, servicesRepo),
		StatuspageBySlugCmd:                 command.NewStatuspageBySlugCmd(lg, servicesRepo),
		IncidentByStatuspageCmd:             command.NewIncidentByStatuspageCmd(lg, servicesRepo, incidentsRepo),
		ScheduledMaintenanceByStatuspageCmd: command.NewScheduledMaintenanceByStatuspageCmd(lg, servicesRepo, scheduledMaintenancesRepo),
	}
	return ServerApplication{
		HttpHandler: handler,
		lg:          lg,
	}
}

func (s ServerApplication) Start(ctx context.Context, addr string) error {
	serverInterface := api.NewStrictHandler(s.HttpHandler, nil)

	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RequestLogger(&CustomLogFormatter{Logger: s.lg}),
		middleware.CleanPath,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
		middleware.Recoverer,
	)

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

type CustomLogFormatter struct {
	Logger *slog.Logger
}

func (l *CustomLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &CustomLogEntry{
		Logger: l.Logger.With(
			"path",
			r.RequestURI,
			"method",
			r.Method,
			"request_id",
			r.Context().Value(middleware.RequestIDKey),
		),
	}

}

type CustomLogEntry struct {
	Logger *slog.Logger
}

func (l *CustomLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	l.Logger.Info(
		"request",
		"status", status,
		"time_taken", fmt.Sprintf("%dms", elapsed.Milliseconds()),
		slog.Any("extra", extra),
	)
}
func (l *CustomLogEntry) Panic(v any, stack []byte) {
	middleware.PrintPrettyStack(v)
}
