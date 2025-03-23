package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"agentos/internal/domain"
)

type Adapter struct {
	*pgxpool.Pool
}

func New(configDB domain.Database) (*Adapter, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		configDB.Driver,
		configDB.User,
		configDB.Password,
		configDB.Host,
		configDB.Port,
		configDB.Name,
		configDB.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("unable to parse config connection: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Adapter{dbPool}, nil
}
