package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey string = os.Getenv("JWT_SECRET_KEY")

type claims struct {
	login string
	exp   int64 `time.Time.Unix()`
	iat   int64 `time.Time.Unix()`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {

		claim := &claims{}
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			prefix := "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				return errors.BadRequest(400, "Wrong auth method", nil, r)
			}
			tokenString := strings.TrimPrefix(authHeader, prefix)
			token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrInvalidKeyType
				}
				return []byte(secretKey), nil
			})
			if err != nil {
				return errors.InternalServerError(500, "Token parse error", err, r)
			}
			if !token.Valid {
				w.Header().Del("Authorization")
				return errors.Unauthorized(401, "Token not valid", nil, r)
			}
			next.ServeHTTP(w, r)
			return nil
		} else {
			return errors.Unauthorized(1001, "Unauthorized", nil, r)
		}
	}))
}

func GetToket(login string) (string, error) {

	claims := claims{
		login: login,
		iat:   time.Now().Unix(),
		exp:   time.Now().Add(4 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
