package scrapper

import (
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type scrapper interface {
	scrap(client *resty.Client) error
}

func New() {
	var providerServices []scrapper
	services, _ := domain.Service.GetAll()

	client := resty.New()

	for _, service := range services {
		if service.ProviderType == types.AtlassianProviderType {
			providerServices = append(providerServices, atlassianProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})
		}
	}

	for _, service := range providerServices {
		if err := service.scrap(client); err != nil {
			slog.Error("error while scraping", err)
		}
	}

}
