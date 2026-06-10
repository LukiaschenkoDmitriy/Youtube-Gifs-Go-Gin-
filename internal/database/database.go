package database

import (
	"context"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(c *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), c.DatabaseUrl)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
