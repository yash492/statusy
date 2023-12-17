package schema

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
