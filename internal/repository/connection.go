package repository

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type ConnectionConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (cfg ConnectionConfig) connectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func Connect(cfg ConnectionConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.connectionString())
	if err != nil {
		return nil, fmt.Errorf("failed connect to the database: %v", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil {
		return nil, err
	}

	// Additional connection settings.
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db, nil
}
