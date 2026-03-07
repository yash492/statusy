package applications

import (
	"log/slog"
	"net/http"

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
	}

}

func (s ServerApplication) Start() {
	serverInterface := api.NewStrictHandler(s.HttpHandler, nil)
	r := chi.NewRouter()
	handler := api.HandlerFromMux(serverInterface, r)
	err := http.ListenAndServe("", handler)
	if err != nil {
		panic(err)
	}
}
