package cmd

type Scrapper interface {
	ScrapIncidents() error
	ScrapMaintenance() error
}
