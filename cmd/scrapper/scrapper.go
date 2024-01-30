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

const workerCount = 10

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

	wg := &sync.WaitGroup{}
	client := resty.New()
	for i, service := range services {
		if i%workerCount == 0 && i != 0 {
			wg.Wait()
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, service scrapper) {
			if err := service.scrap(client, nil); err != nil {
				slog.Error("error while scraping", err)
			}

			wg.Done()
		}(wg, service)
	}

	wg.Wait()

	return nil
}

func scrapStatusPages(queue *queue.Queue, providerServices []scrapper) {
	client := resty.New()

	scrapInterval := config.ScrapIntervalInMins

	ticker := time.NewTicker(time.Duration(scrapInterval) * time.Minute)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			wg := &sync.WaitGroup{}

			for i, service := range providerServices {
				if i%workerCount == 0 && i != 0 {
					wg.Wait()
				}
				wg.Add(1)
				go func(wg *sync.WaitGroup, service scrapper) {
					if err := service.scrap(client, nil); err != nil {
						slog.Error("error while scraping", err)
					}

					wg.Done()
				}(wg, service)
			}
			wg.Wait()
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
		switch service.ProviderType {

		case types.AtlassianProviderType:
			providerServices = append(providerServices, atlassianProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})

		case types.SquadcastProviderType:
			providerServices = append(providerServices, squadcastProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})

		case types.StatusioProviderType:
			providerServices = append(providerServices, statusioProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})
		}

	}

	return providerServices, nil
}
