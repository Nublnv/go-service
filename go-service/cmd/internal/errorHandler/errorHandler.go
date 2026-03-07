package errorHandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	httpErrors "github.com/Nublnv/go-service/cmd/internal/errors"
)

type ApiErrorResponse func(w http.ResponseWriter, r *http.Request) error

func Wrap(h ApiErrorResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var httpErr *httpErrors.HTTPError
			if errors.As(err, &httpErr) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpErr.Status)
				json.NewEncoder(w).Encode(&httpErrors.ErrorResponse{
					Error: httpErrors.ErrorBody{
						Code:    httpErr.Code,
						Message: httpErr.Message,
						Detail:  httpErr.Detail,
					},
				})
				r.Context().Err()
				log.Print(httpErr.Error())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
	}
}
