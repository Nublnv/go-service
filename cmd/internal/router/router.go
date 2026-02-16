package router

import (
	"net/http"

	middleware "github.com/Nublnv/go-service/cmd/internal/middleware"
	public "github.com/Nublnv/go-service/cmd/internal/public"
)

var router *http.ServeMux

func init() {
	router = http.NewServeMux()

	public := public.GetPublicHandler()
	api := http.NewServeMux()

	router.Handle("/", public)
	router.Handle("/api/", middleware.AuthMiddleware(http.StripPrefix("/api", api)))

}

func GetRouter() *http.ServeMux {
	return router
}
