package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ServerBaseContext(ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s: %s %s\n", time.Now().Format("2006-01-02 15:04:05 MST -0700"), r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
