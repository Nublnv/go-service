package api

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/Nublnv/go-service/cmd/internal/validation"
)

var handler http.ServeMux

func init() {
	handler.Handle("/test", errorHandler.Wrap(test))
}

func GetApiHandler() http.Handler {
	return &handler
}

type Test struct {
	Data  string                  `form:"data" validate:"required"`
	File1 *multipart.FileHeader   `file:"file1" validate:"required"`
	File2 []*multipart.FileHeader `file:"test2"`
}

func test(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.MethodNotAllowed(1000, "Method not allowed", nil, r)
	}

	data, err := validation.DecodeAndValidate[Test](r)
	if err != nil {
		return err
	}

	fmt.Println(data.File1.Size)

	return nil
}
