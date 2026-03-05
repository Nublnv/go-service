package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPgxPool(ctx context.Context) (*pgxpool.Pool, error) {

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		panic("POSTGRES_HOST env variable is not set")
	}
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		panic("POSTGRES_USER env variable is not set")
	}
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	if dbPass == "" {
		panic("POSTGRES_PASSWORD env variable is not set")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		panic("POSTGRES_DB env variable is not set")
	}

	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect postgres db with provided credentials")
	}

	return pool, nil
}
