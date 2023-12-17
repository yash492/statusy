package resource

// This package initialises service and components.
// By initialisation, I mean, we insert provider's
// status page service and components
// details to the db for the scrapper and the notification dispatcher
// to work. Since this is both common to dispatcher and scrapper
// I felt a need to create a new package

func New() {
	// Service needs to be initialised  first
	// as components are dependancy of services

	//TODO: error handling
	initServices()
	initComponents()
}
