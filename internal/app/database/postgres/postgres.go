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

func New(cfg config.DatabaseConfig) database.Database {
	return &Postgres{
		cfg: cfg,
	}
}

// TODO: wrap errors
func (p *Postgres) Connect() error {
	db, err := sql.Open("postgres", p.cfg.ConnString)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	p.db = db

	return nil
}
