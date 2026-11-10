package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest[T any](r *http.Request, req *T) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("decode json: %w", err)

	}

	if err := requestValidator.Struct(req); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	return nil
}
