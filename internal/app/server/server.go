package server

import (
	"log"

	"github.com/0x0FACED/tender-service/config"
	"github.com/0x0FACED/tender-service/internal/app/database/postgres"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	servicesimpl "github.com/0x0FACED/tender-service/internal/app/services_impl"
	"github.com/labstack/echo/v4"
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

	bidService := servicesimpl.NewBidService(db)
	tenderService := servicesimpl.NewTenderService(db)
	healthService := &servicesimpl.HealthServiceImpl{}

	s := New(bidService, healthService, tenderService, nil, cfg.Server)

	if err := s.r.Start(s.cfg.Addr); err != nil {
		return err
	}

	return nil
}
