package main

import (
	"github.com/yash492/statusy/cmd/scrapper"
	"github.com/yash492/statusy/pkg/config"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/resource"
	"github.com/yash492/statusy/pkg/store"
)

func main() {
	// This sets up the default logger in the app
	//TODO: error handling and logging
	initLogger()
	config.New()
	store.New()
	domain.New()
	resource.New()
	scrapper.New()

}
