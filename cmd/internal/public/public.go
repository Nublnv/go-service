package public

import (
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
)

var handler http.ServeMux

func init() {
	handler.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	handler.Handle("/register", errorHandler.Wrap(register))
	handler.Handle("/login", errorHandler.Wrap(login))
}

func GetPublicHandler() http.Handler {
	return &handler
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.MethodNotAllowed(1000, "Method not allowed", nil)
	}

	if err := r.ParseForm(); err != nil {
		return errors.BadRequest(1002, "Invalid form data", err)
	}

	registerReq := &RegisterRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var db = r.Context().Value("db")

	if db == nil {
		return errors.InternalServerError(1003, "Database connection not found", nil)
	} else {
		// TODO - implement registration logic here
	}

	w.Header().Set("Authorization", "Bearer <token>") // TODO - generate real token here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successful"))

	return nil
}

func login(w http.ResponseWriter, r *http.Request) error {
	// TODO - implement login logic here
	return nil
}
