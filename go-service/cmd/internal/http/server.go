package http

import (
	"fmt"
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/db"
	"github.com/Nublnv/go-service/cmd/internal/middleware"
	router "github.com/Nublnv/go-service/cmd/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetHttpServer(host string, port string, pool *pgxpool.Pool, chPool *db.Pool) *http.Server {
	svc := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: middleware.DbMiddleware(pool)(middleware.ChMiddleware(chPool)(router.GetRouter())),
	}

	return svc
}

func ServeServer(svc *http.Server, tlsCert string, tlsKey string) func() {
	if tlsCert != "" && tlsKey != "" {
		return func() {
			if err := svc.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("Failed to start server: %v", err))
			}
		}
	} else {
		return func() {
			if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("Failed to start server: %v", err))
			}
		}
	}
}
