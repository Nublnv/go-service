package middleware

import (
	"context"
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
	Login  string `json:"login"`
	Userid int    `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			prefix := "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				return errors.BadRequest(400, "Wrong auth method", nil, r)
			}
			tokenString := strings.TrimPrefix(authHeader, prefix)
			claim, err := checkToken(tokenString, r)
			if err != nil {
				return err
			}
			ctx := context.WithValue(r.Context(), "login", claim.Login)
			ctx = context.WithValue(ctx, "userid", claim.Userid)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		} else {
			return errors.Unauthorized(1001, "Unauthorized", nil, r)
		}
	}))
}

func checkToken(tokenString string, r *http.Request) (*claims, error) {
	claim := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.InternalServerError(500, "Token parse error", err, r)
	}
	if !token.Valid {
		return nil, errors.Unauthorized(401, "Token not valid", nil, r)
	}
	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, errors.Unauthorized(401, "invalid token", nil, r)
	}
	return c, nil
}

func GetToken(login string, userid int, ip string) (string, error) {

	claims := claims{
		Login:  login,
		Userid: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ip,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
