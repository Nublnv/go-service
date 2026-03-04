package api

import (
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
)

var handler http.ServeMux

func init() {
	handler.Handle("/test", errorHandler.Wrap(test))
}

func GetApiHandler() http.Handler {
	return &handler
}

func test(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.MethodNotAllowed(1000, "Method not allowed", nil, r)
	}

	return nil
}
