package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	svc := &http.Server{
		Addr: fmt.Sprintf("%s:%s", host, port),
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Failed to shutdown server: %v", err))
	}

	fmt.Println("Server gracefully stopped")

}
