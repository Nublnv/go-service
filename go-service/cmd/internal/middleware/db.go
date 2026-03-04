package middleware

import (
	"context"
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DbMiddleware(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {
			conn, err := pool.Acquire(r.Context())
			if err != nil {
				return errors.InternalServerError(500, "DB unavailible", err, r)
			}
			defer conn.Release()

			tx, err := conn.Begin(r.Context())
			if err != nil {
				return errors.InternalServerError(500, "Tx begin failed", err, r)
			}
			defer tx.Rollback(r.Context())

			ctx := context.WithValue(r.Context(), "db", tx)
			next.ServeHTTP(w, r.WithContext(ctx))

			if err := tx.Commit(r.Context()); err != nil {
				tx.Rollback(r.Context())
			}
			return nil
		}))
	}
}
