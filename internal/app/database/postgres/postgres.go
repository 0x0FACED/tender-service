package postgres

import (
	"database/sql"

	"github.com/0x0FACED/tender-service/config"
	"github.com/0x0FACED/tender-service/internal/app/database"
	z "github.com/0x0FACED/tender-service/internal/app/logger/zaplog"
	"go.uber.org/zap"
)

type Postgres struct {
	db *sql.DB

	cfg    config.DatabaseConfig
	logger *z.ZapLogger
}

func New(cfg config.DatabaseConfig, logger *z.ZapLogger) database.Database {
	return &Postgres{
		cfg:    cfg,
		logger: logger,
	}
}

// TODO: wrap errors
func (p *Postgres) Connect() error {
	p.logger.Info("Connecting to DB...")
	p.logger.Info("Conn string", zap.String("connstr", p.cfg.ConnString))
	db, err := sql.Open("postgres", p.cfg.ConnString)
	if err != nil {
		p.logger.Error("Error connecting to DB", zap.Error(err))
		return err
	}

	if err := db.Ping(); err != nil {
		p.logger.Error("Error Ping() DB", zap.Error(err))
		return err
	}

	p.db = db

	p.logger.Info("Successfully connected to DB!")
	return nil
}
