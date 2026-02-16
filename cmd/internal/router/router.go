package router

import (
	"net/http"

	middleware "github.com/Nublnv/go-service/cmd/internal/middleware"
)

var router *http.ServeMux

func init() {
	router = http.NewServeMux()

	public := http.NewServeMux()
	public.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	api := http.NewServeMux()

	router.Handle("/", public)
	router.Handle("/api/", middleware.AuthMiddleware(http.StripPrefix("/api", api)))

}

func GetRouter() *http.ServeMux {
	return router
}
