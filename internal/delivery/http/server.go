package http

import (
	"cat-test/internal/usecase"
	"cat-test/internal/validator"
	"context"
	"fmt"
	"net/http"
	"time"
)

type ServerConfig struct {
	Port         uint16        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
}

type server struct {
	srv *http.Server
}

func (s *server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func NewServer(cfg ServerConfig, usecases usecase.Usecases, validators validator.Validators) *server {
	return &server{
		srv: &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			Handler:      getRouter(usecases, validators),
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
		},
	}
}
