package cmd

type Scrapper interface {
	ScrapIncidents()
	ScrapMaintenance()
}
