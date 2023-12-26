package main

import (
	"sync"

	"github.com/yash492/statusy/cmd/dispatcher"
	"github.com/yash492/statusy/cmd/scrapper"
	"github.com/yash492/statusy/pkg/config"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/resource"
	"github.com/yash492/statusy/pkg/store"
)

func main() {

	dispatcherQueue := queue.New(1000)
	wg := sync.WaitGroup{}

	// This sets up the default logger in the app
	//TODO: error handling and logging
	initLogger()
	config.New()
	store.New()
	domain.New()
	resource.New()

	wg.Add(1)
	go scrapper.New(dispatcherQueue, &wg)
	wg.Add(1)
	go dispatcher.New(dispatcherQueue, &wg)
	wg.Wait()

}
