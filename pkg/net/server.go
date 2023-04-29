package net

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/insighted4/correios-cep/pkg/log"
)

const (
	DefaultMaxHeaderBytes  = 1 << 20
	DefaultAddr            = ":8080"
	DefaultReadTimeout     = 20 * time.Second
	DefaultWriteTimeout    = 20 * time.Second
	DefaultIdleTimeout     = 120 * time.Second
	DefaultShutdownTimeout = 5 * time.Minute
)

// HTTPServerConfig holds info required to configure a server.Server.
type HTTPServerConfig struct {
	// MaxHeaderBytes can be used to override the default of 1<<20.
	MaxHeaderBytes int

	// ReadTimeout can be used to override the default http Server timeout of 20s.
	// The string should be formatted like a time.Duration string.
	ReadTimeout time.Duration

	// WriteTimeout can be used to override the default http Server timeout of 20s.
	// The string should be formatted like a time.Duration string.
	WriteTimeout time.Duration

	// IdleTimeout can be used to override the default http Server timeout of 120s.
	// The string should be formatted like a time.Duration string.
	IdleTimeout time.Duration

	// ShutdownTimeout can be used to override the default http Server Shutdown timeout
	// of 5m.
	ShutdownTimeout time.Duration

	// Addr is the binding address the Server implementation will serve HTTP over.
	// The default is ":8080"
	Addr string
}

type Server interface {
	Run() error
}

// server encapsulates all logic for registering and running a Server.
type server struct {
	cfg HTTPServerConfig

	httpServer *http.Server
	Shutdown   func()

	// Exit chan for graceful Shutdown
	Exit chan chan error
}

func NewServer(httpCfg HTTPServerConfig, handler http.Handler, shutdownFn func()) *server {
	if httpCfg.MaxHeaderBytes == 0 {
		httpCfg.MaxHeaderBytes = DefaultMaxHeaderBytes
	}

	if httpCfg.ReadTimeout == 0 {
		httpCfg.ReadTimeout = DefaultReadTimeout
	}

	if httpCfg.WriteTimeout == 0 {
		httpCfg.WriteTimeout = DefaultWriteTimeout
	}

	if httpCfg.IdleTimeout == 0 {
		httpCfg.IdleTimeout = DefaultIdleTimeout
	}

	if httpCfg.ShutdownTimeout == 0 {
		httpCfg.ShutdownTimeout = DefaultShutdownTimeout
	}

	if httpCfg.Addr == "" {
		httpCfg.Addr = DefaultAddr
	}

	httpServer := &http.Server{
		Handler:        handler,
		Addr:           httpCfg.Addr,
		MaxHeaderBytes: httpCfg.MaxHeaderBytes,
		ReadTimeout:    httpCfg.ReadTimeout,
		WriteTimeout:   httpCfg.WriteTimeout,
		IdleTimeout:    httpCfg.IdleTimeout,
	}

	return &server{
		cfg:        httpCfg,
		httpServer: httpServer,
		Shutdown:   shutdownFn,
		Exit:       make(chan chan error),
	}
}

func (s *server) start() error {
	go func() {
		log.Infof("Listening and serving HTTP on %s", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP Server error - initiating shutting down: %v", err)
			s.stop()
			return
		}
	}()

	go func() {
		exit := <-s.Exit

		// Stop listener with timeout
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		// Stop service
		if s.Shutdown != nil {
			s.Shutdown()
		}

		// Stop HTTP Server
		if s.httpServer != nil {
			log.Infof("Stopping HTTP Server on %s", s.httpServer.Addr)
			exit <- s.httpServer.Shutdown(ctx)
			return
		}

		exit <- nil
	}()

	return nil
}

func (s *server) stop() error {
	ch := make(chan error)
	s.Exit <- ch
	return <-ch
}

// Run will create a new Server and register the given
// Service and start up the Server(s).
// This will block until the Server shuts down.
func (s *server) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Received signal ", <-ch)
	return s.stop()
}
