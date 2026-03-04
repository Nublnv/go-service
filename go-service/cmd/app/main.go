package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	host := os.Getenv("REST_HOST")
	if host == "" {
		panic("REST_HOST env variable is not set")
	}
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "443"
	}

	tlsPath := os.Getenv("TLS_PATH")
	if tlsPath == "" {
		panic("tlsPath env variables must be set with the path that include server.crt and server.key")
	} else {
		stat, err := os.Stat(tlsPath)
		if err != nil {
			panic(err)
		}
		if !stat.IsDir() {
			panic("TLS_PATH must be a directory")
		}
		if _, err := os.Stat(fmt.Sprintf("%s/server.crt", tlsPath)); err != nil {
			panic(err)
		}
		if _, err := os.Stat(fmt.Sprintf("%s/server.key", tlsPath)); err != nil {
			panic(err)
		}
	}

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

	baseContext := context.Background()

	pool, err := pgxpool.New(baseContext, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		panic("Failed to connect postgres db with provided credentials")
	}

	svc := http.GetHttpServer(host, port, pool)
	go http.ServeServer(svc, fmt.Sprintf("%s/server.crt", tlsPath), fmt.Sprintf("%s/server.key", tlsPath))()

	fmt.Printf("Server is running on %s:%s\n", host, port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(baseContext, 5*time.Second)
	defer pool.Close()
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Failed to shutdown server: %v", err))
	}

	fmt.Println("Server gracefully stopped")

}
