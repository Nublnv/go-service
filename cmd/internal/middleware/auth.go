package middleware

import (
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// TODO - implement auth logic here
			next.ServeHTTP(w, r)
			return nil
		} else {
			return errors.Unauthorized(1001, "Unauthorized", nil)
		}
	}))
}
