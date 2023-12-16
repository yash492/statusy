package store

type serviceDBConn struct {
	db
}

func NewServiceDBConn() serviceDBConn {
	return serviceDBConn{
		db: dbConn,
	}
}
