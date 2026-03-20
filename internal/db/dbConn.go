package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(DBUrl string) (*pgxpool.Pool, error) {
   config, err := pgxpool.ParseConfig(DBUrl)
   if err != nil {
	slog.Error("Unable to parse database URL", "error", err)
	return nil, err
   }

   config.MaxConns = 10
   config.MinConns = 2
   config.MaxConnLifetime = 0
   config.MaxConnIdleTime = 5 * time.Second
   
   ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
   defer cancel()
   pool, err := pgxpool.NewWithConfig(ctx, config)
   if err != nil {
	slog.Error("Unable to create connection pool", "error", err)
	return nil, err
   }

   return pool, nil
}