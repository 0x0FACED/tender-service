package postgres

import (
	"database/sql"

	"github.com/0x0FACED/tender-service/config"
	"github.com/0x0FACED/tender-service/internal/app/database"
)

type Postgres struct {
	db *sql.DB

	cfg config.DatabaseConfig
}

func New() database.Database {
	// ...
	return nil
}
