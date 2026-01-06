package depositor

import "github.com/yash492/statusy/internal/domain/statuspage"

type Depositor interface {
	SaveComponents() (statuspage.AggregateComponents, error)
	SaveIncidents() ([]statuspage.Incident, error)
}
