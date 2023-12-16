package store

type componentDBConn struct {
	db
}

func NewComponentDBConn() componentDBConn {
	return componentDBConn{
		db: dbConn,
	}
}
