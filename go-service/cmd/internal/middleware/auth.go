package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	httpErorrs "github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var secretKey string = os.Getenv("JWT_SECRET_KEY")

type claims struct {
	Login  string `json:"login"`
	Userid int    `json:"user_id"`
	jwt.RegisteredClaims
}

type parseError struct {
	err error
}
type invalidError struct {
	err error
}

func (e *parseError) Http(r *http.Request) error {
	return httpErorrs.InternalServerError(500, "Token parse error", e.err, r)
}
func (*parseError) Rpc() error {
	return status.Error(codes.Internal, "Token parse error")
}
func (e *parseError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "Token parse error"
}
func (e *parseError) Unwrap() error {
	return e.err
}
func (e *invalidError) Http(r *http.Request) error {
	return httpErorrs.InternalServerError(401, "Invalid token", e.err, r)
}
func (*invalidError) Rpc() error {
	return status.Error(codes.Unauthenticated, "Invalid token")
}
func (e *invalidError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "Invalid token"
}
func (e *invalidError) Unwrap() error {
	return e.err
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			prefix := "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				return httpErorrs.BadRequest(400, "Wrong auth method", nil, r)
			}
			tokenString := strings.TrimPrefix(authHeader, prefix)
			claim, err := checkToken(tokenString)
			if err != nil {
				var parseErr *parseError
				var invalidErr *invalidError
				if errors.As(err, &parseErr) {
					return parseErr.Http(r)
				}
				if errors.As(err, &invalidErr) {
					return invalidErr.Http(r)
				}
			}
			ctx := context.WithValue(r.Context(), "login", claim.Login)
			ctx = context.WithValue(ctx, "userid", claim.Userid)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		} else {
			return httpErorrs.Unauthorized(1001, "Unauthorized", nil, r)
		}
	}))
}

func checkToken(tokenString string) (*claims, error) {
	claim := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, &parseError{err: err}
	}
	if !token.Valid {
		return nil, &invalidError{}
	}
	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, &invalidError{}
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

func AuthMiddlewareRpc(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}
	authHeaderArr := md.Get("authorization")
	if len(authHeaderArr) == 1 {
		authHeader := authHeaderArr[0]
		prefix := "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			return nil, status.Error(codes.Unauthenticated, "Wrong ")
		}
		tokenString := strings.TrimPrefix(authHeader, prefix)
		claim, err := checkToken(tokenString)
		if err != nil {
			var parseErr *parseError
			var invalidErr *invalidError
			if errors.As(err, &parseErr) {
				return nil, parseErr.Rpc()
			}
			if errors.As(err, &invalidErr) {
				return nil, invalidErr.Rpc()
			}
		}
		ctx = context.WithValue(ctx, "login", claim.Login)
		ctx = context.WithValue(ctx, "userid", claim.Userid)
		return handler(ctx, req)
	} else {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}

}
