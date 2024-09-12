package server

import (
	"github.com/0x0FACED/tender-service/config"
	"github.com/0x0FACED/tender-service/internal/app/database/postgres"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/0x0FACED/tender-service/internal/app/logger/zaplog"
	servicesimpl "github.com/0x0FACED/tender-service/internal/app/services_impl"
	"github.com/0x0FACED/tender-service/migrations"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type server struct {
	r *echo.Echo

	bidHandler    repos.BidService
	healthHandler repos.HealthService
	tenderHandler repos.TenderService

	logger *zaplog.ZapLogger
	cfg    config.ServerConfig
}

func New(
	bid repos.BidService,
	health repos.HealthService,
	tender repos.TenderService,
	logger *zaplog.ZapLogger,
	cfg config.ServerConfig,

) *server {
	return &server{
		r:             echo.New(),
		bidHandler:    bid,
		healthHandler: health,
		tenderHandler: tender,
		logger:        logger,
		cfg:           cfg,
	}
}

func Start() error {
	l := zaplog.New()
	l.Info("Starting server...")
	cfg, err := config.Load()
	if err != nil {
		l.Fatal("cant load config", zap.Error(err))
		return err
	}

	l.Info("Config successfully loaded")

	// Можно добавить в дальнейшем выбр базы данных исходя из .env
	db := postgres.New(cfg.Database, l)
	if err := db.Connect(); err != nil {
		l.Fatal("cant connect to db", zap.Error(err))
		return err
	}

	l.Info("DB successfully connected")

	bidService := servicesimpl.NewBidService(db)
	tenderService := servicesimpl.NewTenderService(db)
	healthService := &servicesimpl.HealthServiceImpl{}

	if err := migrations.Up(cfg.Database.ConnString); err != nil {
		l.Fatal("cant migrate up", zap.Error(err))
		return err
	}

	l.Info("Migrate Up successfully")

	s := New(bidService, healthService, tenderService, l, cfg.Server)
	s.RegisterHandlers()
	s.r.Use(middleware.Logger())

	l.Info("Server created, handlers registered, using middleware: middleware.Logger()")

	l.Info("Starting listen on addr", zap.String("addr", s.cfg.Addr))
	if err := s.r.Start(s.cfg.Addr); err != nil {
		l.Fatal("Server error", zap.Error(err))
		return err
	}

	return nil
}
