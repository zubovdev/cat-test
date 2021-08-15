package main

import (
	"cat-test/internal/delivery/http"
	"cat-test/internal/reminder"
	"cat-test/internal/repository"
	"cat-test/internal/usecase"
	"cat-test/internal/validator"
	"context"
	"crypto/tls"
	"github.com/caarlos0/env/v6"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	HTTPPort         uint16        `env:"HTTP_PORT" envDefault:"8080"`
	HTTPReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	HTTPWriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`

	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     uint16 `env:"DB_PORT" envDefault:"5435"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"secret"`
	DBDatabase string `env:"DB_DATABASE" envDefault:"postgres"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env file: %v", err)
	}

	dialer := gomail.NewDialer("in-v3.mailjet.com", 587, "72585d717626390fb1e0e27cdc83ce18", "619de4d1868c37fb8fda889bb6211e23")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	db, err := repository.Connect(repository.ConnectionConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBDatabase,
	})
	if err != nil {
		log.Fatal(err)
	}


	repositories := repository.NewRepositories(db)
	validators := validator.NewValidators(db)
	usecases := usecase.NewUsecases(repositories)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go reminder.RunTaskReminder(ctx, usecases.Task, dialer)

	srv := http.NewServer(http.ServerConfig{
		Port:         cfg.HTTPPort,
		WriteTimeout: cfg.HTTPWriteTimeout,
		ReadTimeout:  cfg.HTTPReadTimeout,
	}, usecases, validators)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown http server: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}

	os.Exit(0)
}
