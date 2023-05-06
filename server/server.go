package server

import (
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/insighted4/correios-cep/correios"
	"github.com/insighted4/correios-cep/pkg/app"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/health"
	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/pkg/net"
	"github.com/insighted4/correios-cep/pkg/version"
	"github.com/insighted4/correios-cep/server/handler"
	"github.com/insighted4/correios-cep/storage"
	"github.com/sirupsen/logrus"
)

type Server interface {
	Run() error
	Shutdown()
}

type Config struct {
	HTTPServerConfig net.HTTPServerConfig

	// Set gin mode to release.
	ReleaseMode bool

	Storage storage.Storage

	// If specified, the server will use this function for determining time.
	Now func() time.Time
}

type Service struct {
	cfg      Config
	correios correios.Correios
	health   gosundheit.Health
	logger   logrus.FieldLogger
	server   net.Server
	storage  storage.Storage

	now func() time.Time
}

var _ Server = (*Service)(nil)

func New(cfg Config) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}

	correios := correios.New(net.NewClient())
	healthChecker := gosundheit.New()

	svc := &Service{
		cfg:      cfg,
		correios: correios,
		health:   healthChecker,
		logger:   log.WithField("component", "server"),
		storage:  cfg.Storage,
		now:      cfg.Now,
	}

	httpHandler := handler.New(correios, cfg.Storage, healthChecker, cfg.ReleaseMode, cfg.Now)

	svc.server = net.NewServer(cfg.HTTPServerConfig, httpHandler, svc.Shutdown)

	return svc
}

func (s *Service) Run() error {
	const op errors.Op = "server.Run"
	s.logger.Infof("%s: Starting HTTP Server (%s)", app.Description, version.Version)

	if err := s.health.RegisterCheck(&checks.CustomCheck{
		CheckName: "correios",
		CheckFunc: health.NewCustomHealthCheckFunc(s.correios, s.now),
	}, gosundheit.ExecutionPeriod(1*time.Second),
		gosundheit.InitiallyPassing(false)); err != nil {
		return errors.E(op, errors.KindUnexpected, err)
	}

	if s.storage != nil {
		if err := s.health.RegisterCheck(&checks.CustomCheck{
			CheckName: "database",
			CheckFunc: health.NewCustomHealthCheckFunc(s.storage, s.now),
		}, gosundheit.ExecutionPeriod(10*time.Second),
			gosundheit.InitiallyPassing(false)); err != nil {
			return errors.E(op, errors.KindUnexpected, err)
		}
	} else {
		return errors.E(op, errors.KindUnexpected, "invalid storage configuration")
	}

	// Start Server
	if err := s.server.Run(); err != nil {
		return errors.E(op, "failed to start server", err)
	}

	return nil
}

func (s *Service) Shutdown() {
	s.logger.Infof("%s: Stopping HTTP Server", app.Description)
	s.storage.Close()
}
