package server

import (
	"log"

	"github.com/0x0FACED/tender-service/config"
	"github.com/0x0FACED/tender-service/internal/app/database/postgres"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
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

	logger *zap.Logger
	cfg    config.ServerConfig
}

func New(
	bid repos.BidService,
	health repos.HealthService,
	tender repos.TenderService,
	logger *zap.Logger,
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Cant load config file, return")
		return err
	}

	db := postgres.New(cfg.Database)
	if err := db.Connect(); err != nil {
		log.Fatalln("Cant connect to DB, err: ", err)
	}

	bidService := servicesimpl.NewBidService(db)
	tenderService := servicesimpl.NewTenderService(db)
	healthService := &servicesimpl.HealthServiceImpl{}

	if err := migrations.Up(cfg.Database.ConnString); err != nil {
		log.Fatalln("cant migrate up, err: ", err)
	}

	s := New(bidService, healthService, tenderService, nil, cfg.Server)

	s.RegisterHandlers()

	s.r.Use(middleware.Logger())

	if err := s.r.Start(s.cfg.Addr); err != nil {
		return err
	}

	return nil
}
