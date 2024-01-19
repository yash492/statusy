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

func New(queue *queue.Queue, wg *sync.WaitGroup) {
	var providerServices []scrapper
	services, _ := domain.Service.GetAll()

	for _, service := range services {
		if service.ProviderType == types.AtlassianProviderType {
			providerServices = append(providerServices, atlassianProvider{
				incidentUrl: service.IncidentURL,
				serviceID:   service.ID,
			})
		}
	}

	scrapStatusPages(queue, providerServices)
	wg.Done()

}

func scrapStatusPages(queue *queue.Queue, providerServices []scrapper) {
	client := resty.New()
	client.SetTimeout(time.Duration(1 * time.Minute))

	scrapInterval := config.ScrapIntervalInMins

	ticker := time.NewTicker(time.Duration(scrapInterval) * time.Second)
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
