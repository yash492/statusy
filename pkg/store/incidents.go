package store

type incidentDBConn struct {
	db
}

func NewIncidentDBConn() incidentDBConn {
	return incidentDBConn{
		db: dbConn,
	}
}

func (db incidentDBConn) Try() {

}
