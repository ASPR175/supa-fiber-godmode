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

	// Custom dialer that forces IPv4
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			// Force only IPv4 lookups
			// This works because it avoids AAAA queries entirely
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("tcp4", address)
			},
		},
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	// Override the dialer
	config.ConnConfig.DialFunc = dialer.DialContext

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	DB = pool
	return pool
}
