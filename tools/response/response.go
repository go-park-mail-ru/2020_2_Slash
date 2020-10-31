package response

import "github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"

type Body map[string]interface{}

type Response struct {
	Error   *errors.Error `json:"error,omitempty"`
	Message string        `json:"message,omitempty"`
	Body    *Body         `json:"body,omitempty"`
}
