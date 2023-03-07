package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ethpandaops/stubbies/pkg/api"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log *logrus.Logger
	Cfg Config

	http *api.Handler
}

func NewServer(log *logrus.Logger, conf *Config) *Server {
	if err := conf.Validate(); err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	s := &Server{
		Cfg: *conf,
		log: log,

		http: api.NewHandler(log, &conf.Execution),
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Infof("starting stubbies server")

	s.http.Start(ctx)

	router := httprouter.New()

	if err := s.http.Register(ctx, router); err != nil {
		return err
	}

	if err := s.ServeMetrics(ctx); err != nil {
		return err
	}

	server := &http.Server{
		Addr:              s.Cfg.Addr,
		ReadHeaderTimeout: 3 * time.Minute,
		WriteTimeout:      15 * time.Minute,
	}

	server.Handler = router

	s.log.Infof("serving http at %s", s.Cfg.Addr)

	if err := server.ListenAndServe(); err != nil {
		s.log.Fatal(err)
	}

	return nil
}

func (s *Server) ServeMetrics(ctx context.Context) error {
	go func() {
		server := &http.Server{
			Addr:              s.Cfg.MetricsAddr,
			ReadHeaderTimeout: 15 * time.Second,
		}

		server.Handler = promhttp.Handler()

		s.log.Infof("serving metrics at %s", s.Cfg.MetricsAddr)

		if err := server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}()

	return nil
}
