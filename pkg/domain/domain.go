package domain

import "github.com/yash492/statusy/pkg/store"

var Incident store.IncidentStore
var Component store.ComponentStore
var Service store.ServiceStore

func New() {
	Incident = store.NewIncidentDBConn()
	Component = store.InitDbVar()
	Service = store.InitDbVar()
}
