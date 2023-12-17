package main

import (
	"github.com/yash492/statusy/pkg/config"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/resource"
	"github.com/yash492/statusy/pkg/store"
)

func main() {
	// This sets up the default logger in the app
	setupSlog()
	config.New()
	store.New()
	domain.New()

	resource.InitServices()

}
