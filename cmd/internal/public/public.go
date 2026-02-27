package public

import (
	"context"
	"net/http"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

	var db pgxpool.Tx = r.Context().Value("db").(pgxpool.Tx)

	var existedLogin string = ""

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res := db.QueryRow(ctx, "SELECT login FROM auth.users WHERE login = %s", registerReq.Username)
	res.Scan(&existedLogin)
	if existedLogin != "" {
		return errors.BadRequest(400, "User with provided login already exists", nil)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError(500, "Cannot hash password", err)
	}

	insert, err := db.Exec(ctx, "INSERT INTO auth.users (login, passhash) VALUES ( %s, %d )", registerReq.Username, passHash)
	if err != nil || insert.RowsAffected() == 0 {
		return errors.InternalServerError(500, "Cannot add new user", err)
	}

	return nil
}

func login(w http.ResponseWriter, r *http.Request) error {
	// TODO - implement login logic here
	return nil
}
