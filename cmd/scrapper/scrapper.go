package scrapper

import (
	"log/slog"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yash492/statusy/pkg/config"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/types"
)

type scrapper interface {
	scrap(client *resty.Client, queue *queue.Queue) error
}

func New(queue *queue.Queue, wg *sync.WaitGroup) error {
	services, err := getServices()
	if err != nil {
		return err
	}
	scrapStatusPages(queue, services)
	wg.Done()
	return nil

}

func ScrapStatusPagesDuringAppInitialization() error {
	services, err := getServices()
	if err != nil {
		return err
	}
	client := resty.New()
	client.SetTimeout(time.Duration(1 * time.Minute))

	for _, service := range services {
		if err := service.scrap(client, nil); err != nil {
			slog.Error("error while scraping", err)
		}
	}

	return nil
}

func scrapStatusPages(queue *queue.Queue, providerServices []scrapper) {
	client := resty.New()
	client.SetTimeout(time.Duration(1 * time.Minute))

	scrapInterval := config.ScrapIntervalInMins

	ticker := time.NewTicker(time.Duration(scrapInterval) * time.Minute)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			for _, service := range providerServices {
				if err := service.scrap(client, queue); err != nil {
					slog.Error("error while scraping", err)
				}
			}
		}
	}
}

func getServices() ([]scrapper, error) {
	var providerServices []scrapper
	services, err := domain.Service.GetAll()
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if service.ProviderType == types.AtlassianProviderType {
			providerServices = append(providerServices, atlassianProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})
		}

		if service.ProviderType == types.SquadcastProviderType {
			providerServices = append(providerServices, squadcastProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})
		}
	}

	return providerServices, nil
}
