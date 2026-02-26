package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/middleware"
	router "github.com/Nublnv/go-service/cmd/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	host := os.Getenv("REST_HOST")
	if host == "" {
		panic("REST_HOST env variable is not set")
	}
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8443"
	}

	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")
	if (tlsCert == "" && tlsKey != "") || (tlsCert != "" && tlsKey == "") {
		panic("TLS_CERT and TLS_KEY env variables must be set both or not set at all")
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
	dbPass := os.Getenv("POSTGRES_PASS")
	if dbPass == "" {
		panic("POSTGRES_PASS env variable is not set")
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

	svc := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: middleware.DbMiddleware(pool)(router.GetRouter()),
	}

	if tlsCert != "" && tlsKey != "" {
		go func() {
			if err := svc.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("Failed to start server: %v", err))
			}
		}()
	} else {
		go func() {
			if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("Failed to start server: %v", err))
			}
		}()
	}

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
