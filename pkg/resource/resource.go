package resource

import "github.com/yash492/statusy/cmd/scrapper"

// This package initialises service and components.
// By initialisation, I mean, we insert provider's
// status page service and components
// details to the db for the scrapper and the notification dispatcher
// to work. Since this is both common to dispatcher and scrapper
// I felt a need to create a new package

func New() error {
	// Service needs to be initialised  first
	// as components are dependancy of services

	//TODO: error handling
	err := initServices()
	if err != nil {
		return err
	}
	
	err = initComponents()
	if err != nil {
		return err
	}

	err = scrapper.ScrapStatusPagesDuringAppInitialization()
	if err != nil {
		return err
	}

	return nil
}
