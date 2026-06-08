package applications

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy"
	"github.com/yash492/statusy/internal/adapter/pgx/componentgroupsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/componentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/notificationsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/scheduledmaintenancesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/viewsdb"
	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/common/apperrors"
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
	viewsRepo := viewsdb.NewPostgresViewsRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	componentsRepo := componentsdb.NewPostgresComponentRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	componentGroupsRepo := componentgroupsdb.NewPostgresComponentGroupsRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	notificationsRepo := notificationsdb.NewPostgresNotificationsRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	handler := httphandler.Handler{
		ListStatuspageCmd:                   command.NewListStatuspageCmd(lg, servicesRepo),
		StatuspageBySlugCmd:                 command.NewStatuspageBySlugCmd(lg, servicesRepo),
		IncidentByStatuspageCmd:             command.NewIncidentByStatuspageCmd(lg, servicesRepo, incidentsRepo),
		ScheduledMaintenanceByStatuspageCmd: command.NewScheduledMaintenanceByStatuspageCmd(lg, servicesRepo, scheduledMaintenancesRepo),
		GetOrCreateDefaultViewCmd:           command.NewGetOrCreateDefaultViewCmd(lg, viewsRepo),
		GetUnconfiguredServicesCmd:          command.NewGetUnconfiguredServicesCmd(lg, viewsRepo),
		GetServiceComponentsCmd:             command.NewGetServiceComponentsCmd(lg, servicesRepo, componentsRepo, componentGroupsRepo),
		AddViewServiceCmd:                   command.NewAddViewServiceCmd(lg, viewsRepo),
		EditViewServiceCmd:                  command.NewEditViewServiceCmd(lg, viewsRepo),
		GetViewServiceCmd:                   command.NewGetViewServiceCmd(lg, viewsRepo),
		DeleteViewServiceCmd:                command.NewDeleteViewServiceCmd(lg, viewsRepo),
		EditViewCmd:                         command.NewEditViewCmd(lg, viewsRepo),
		DeleteViewCmd:                       command.NewDeleteViewCmd(lg, viewsRepo),
		GetViewServicesCmd:                  command.NewGetViewServicesCmd(lg, viewsRepo),
		ListViewsCmd:                        command.NewListViewsCmd(lg, viewsRepo),
		CreateViewCmd:                       command.NewCreateViewCmd(lg, viewsRepo),
		GetViewCmd:                          command.NewGetViewCmd(lg, viewsRepo),
		AddViewNotificationCmd:              command.NewAddViewNotificationHandler(lg, notificationsRepo, viewsRepo),
		GetViewNotificationsCmd:             command.NewGetViewNotificationsHandler(lg, notificationsRepo, viewsRepo),
		EditViewNotificationCmd:             command.NewEditViewNotificationHandler(lg, notificationsRepo),
		DeleteViewNotificationCmd:           command.NewDeleteViewNotificationHandler(lg, notificationsRepo),
	}
	return ServerApplication{
		HttpHandler: handler,
		lg:          lg,
	}
}

func (s ServerApplication) Start(ctx context.Context, addr string) error {
	options := api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{
				Code:    string(apperrors.TypeInvalidInput),
				Message: err.Error(),
			})
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			if appErr, isError := errors.AsType[*apperrors.AppError](err); isError {
				s.lg.ErrorContext(r.Context(), "request failed",
					slog.String("type", string(appErr.Type)),
					slog.Any("err", err),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(appErr.StatusCode)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{
					Code:    string(appErr.Type),
					Message: appErr.Message,
				})
				return
			}

			// Unhandled error — safe 500
			s.lg.ErrorContext(r.Context(), "unhandled error", slog.Any("err", err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{
				Code:    string(apperrors.TypeInternal),
				Message: "internal server error",
			})
		},
	}
	serverInterface := api.NewStrictHandlerWithOptions(s.HttpHandler, nil, options)

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

	r.Get("/*", serveStatic(statusy.FrontendFs, "_ui/build"))

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

func serveStatic(embedFS embed.FS, directory string) http.HandlerFunc {
	subFS, err := fs.Sub(embedFS, directory)
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(subFS))

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)
		filePath := strings.TrimPrefix(path, "/")

		if filePath == "" {
			filePath = "index.html"
		}

		file, err := subFS.Open(filePath)
		if err != nil {
			// File does not exist, serve index.html for SPA router fallback
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.Copy(w, indexFile)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			// If path points to a directory, fallback to index.html
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.Copy(w, indexFile)
			return
		}

		fileServer.ServeHTTP(w, r)
	}
}
