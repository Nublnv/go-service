package validation

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	httpErrors "github.com/Nublnv/go-service/cmd/internal/errors"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// Decode and validate data from body/form/form-data/query. Returning T and *HTTPError
func DecodeAndValidate[T any](r *http.Request) (*T, error) {
	ct := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return decodeAndValidateBody[T](r)
	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		return decodeAndValidateForm[T](r)
	case strings.HasPrefix(ct, "multipart/form-data"):
		return decodeAndValidateMultipartFormData[T](r)
	default:
		return decodeAndValidateQuery[T](r)
	}
}

// Validate and decode body, returning T and *httpErrors.HTTPError
func decodeAndValidateBody[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	var dst T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, httpErrors.ValidationError(422, "Body is empty", err, r)
		}

		return nil, httpErrors.ValidationError(422, "Something went wrong", err, r)
	}

	if decoder.More() {
		return nil, httpErrors.ValidationError(422, "Body must contains only one object", nil, r)
	}
	if err := validate(&dst, r); err != nil {
		return nil, err
	}

	return &dst, nil
}

func decodeAndValidateForm[T any](r *http.Request) (*T, error) {

	decoder := form.NewDecoder()

	if err := r.ParseForm(); err != nil {
		return nil, httpErrors.ValidationError(422, "Cannot parse form", err, r)
	}

	var dst T

	if err := decoder.Decode(&dst, r.Form); err != nil {
		return nil, httpErrors.ValidationError(422, "Cannot parse form", err, r)
	}
	if err := validate(&dst, r); err != nil {
		return nil, err
	}

	return &dst, nil
}

func decodeAndValidateQuery[T any](r *http.Request) (*T, error) {
	decoder := form.NewDecoder()

	var dst T

	if err := decoder.Decode(&dst, r.URL.Query()); err != nil {
		return nil, httpErrors.ValidationError(422, "Cannot parse query params", err, r)
	}
	if err := validate(&dst, r); err != nil {
		return nil, err
	}

	return &dst, nil
}

func decodeAndValidateMultipartFormData[T any](r *http.Request) (*T, error) {
	decoder := form.NewDecoder()

	var dst T

	if err := r.ParseMultipartForm(50); err != nil {
		return nil, httpErrors.InternalServerError(500, "Cannot parse form data", err, r)
	}

	if err := decoder.Decode(&dst, r.MultipartForm.Value); err != nil {
		return nil, httpErrors.ValidationError(422, "Cannot parse form data", err, r)
	}

	addFilesFromForm(&dst, r)

	if err := validate(&dst, r); err != nil {
		return nil, err
	}

	return &dst, nil
}

func addFilesFromForm(dst any, r *http.Request) {
	t := reflect.TypeOf(dst)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(dst).Elem()

	var fileHeaderType = reflect.TypeOf((*multipart.FileHeader)(nil))

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		filename := field.Tag.Get("file")
		if filename == "" {
			continue
		}
		val, ok := r.MultipartForm.File[filename]
		if !ok {
			continue
		}
		if field.Type == fileHeaderType {
			valueField := v.FieldByIndex(field.Index)
			valueField.Set(reflect.ValueOf(val[0]))
		}
		if field.Type.Kind() == reflect.Slice && field.Type.Elem() == fileHeaderType {
			valueField := v.FieldByIndex(field.Index)
			valueField.Set(reflect.ValueOf(val))
		}
	}
}

func validate(dst any, r *http.Request) error {
	var validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string
		if field := fld.Tag.Get("json"); field != "" {
			name = strings.Split(field, ",")[0]
		} else if field := fld.Tag.Get("form"); field != "" {
			name = strings.Split(field, ",")[0]
		} else if field := fld.Tag.Get("file"); field != "" {
			name = strings.Split(field, ",")[0]
		}
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(dst); err != nil {
		verr, ok := err.(validator.ValidationErrors)
		if !ok {
			return httpErrors.ValidationError(422, verr.Error(), verr, r)
		}
		fields := make(map[string][]string)
		for _, field := range verr {
			switch field.Tag() {
			case "required":
				fields["field required"] = append(fields["field required"], field.Field())
			default:
				fields["is invalid"] = append(fields["is invalid"], field.Field())
			}
		}
		return httpErrors.ValidationError(422, "Invalid data", nil, r, validationMessage(fields)...)
	}
	return nil
}

func validationMessage(fields map[string][]string) []httpErrors.ValidationErrorDetail {
	result := []httpErrors.ValidationErrorDetail{}

	for k, v := range fields {
		result = append(result, httpErrors.ValidationErrorDetail{
			Fields: v,
			Msg:    k,
		})
	}
	return result
}
