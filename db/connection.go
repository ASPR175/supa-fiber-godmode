package db

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,

			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("tcp4", address)
			},
		},
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	config.ConnConfig.DialFunc = dialer.DialContext

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	DB = pool
	return pool
}
