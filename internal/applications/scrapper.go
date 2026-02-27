package applications

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Logger  *slog.Logger
	ReadDB  *pgxpool.Pool
	WriteDB *pgxpool.Pool
}

func StartScrapper(deps Deps) {

	// incidentResult :=

}

func LoadServicesFromYaml(ctx context.Context, lg *slog.Logger) ([]byte, error) {
	filePath := "../../data/services.yaml"
	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		lg.ErrorContext(ctx, "error while reading sevices from yaml", slog.Any("err", err))
		return nil, err
	}

	return (yamlBytes), nil
}
