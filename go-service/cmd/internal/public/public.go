package public

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/Nublnv/go-service/cmd/internal/middleware"
	"github.com/jackc/pgx/v5"
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
	Username string `json:"login"`
	Password string `json:"password"`
}

type UserData struct {
	login    string
	passhash []byte
}

func register(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.MethodNotAllowed(1000, "Method not allowed", nil, r)
	}

	if err := r.ParseForm(); err != nil {
		return errors.BadRequest(1002, "Invalid form data", err, r)
	}

	registerReq := &RegisterRequest{
		Username: r.FormValue("login"),
		Password: r.FormValue("password"),
	}

	var db pgx.Tx = r.Context().Value("pg").(pgx.Tx)

	var existedLogin string = ""

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res := db.QueryRow(ctx, "SELECT login FROM auth.users WHERE login = $1", registerReq.Username)
	err := res.Scan(&existedLogin)
	if err != nil && err != pgx.ErrNoRows {
		return errors.BadRequest(500, "Something went wrong", err, r)
	}
	if existedLogin != "" {
		return errors.BadRequest(400, "User with provided login already exists", nil, r)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError(500, "Cannot hash password", err, r)
	}

	insert, err := db.Exec(ctx, "INSERT INTO auth.users (login, passhash) VALUES ( $1, $2 )", registerReq.Username, passHash)
	if err != nil || insert.RowsAffected() == 0 {
		return errors.InternalServerError(500, "Cannot add new user", err, r)
	}

	return nil
}

func login(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		errors.MethodNotAllowed(1000, "Method not allowed", nil, r)
	}

	if err := r.ParseForm(); err != nil {
		return errors.BadRequest(1002, "Invalid form data", err, r)
	}

	loginReq := &RegisterRequest{
		Username: r.FormValue("login"),
		Password: r.FormValue("password"),
	}

	var db pgx.Tx = r.Context().Value("pg").(pgx.Tx)

	userData := &UserData{}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res := db.QueryRow(ctx, "SELECT login, passhash FROM auth.users WHERE login = $1", loginReq.Username)
	err := res.Scan(&userData.login, &userData.passhash)
	if err != nil && err != pgx.ErrNoRows {
		return errors.BadRequest(500, "Something went wrong", err, r)
	} else if err == pgx.ErrNoRows {
		return errors.BadRequest(400, "User not found", nil, r)
	}

	err = bcrypt.CompareHashAndPassword(userData.passhash, []byte(loginReq.Password))
	if err != nil {
		return errors.Forbidden(403, "Wrong password", err, r)
	}

	token, err := middleware.GetToken(userData.login)
	if err != nil {
		return errors.InternalServerError(500, "Internal server error", err, r)
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil
}
