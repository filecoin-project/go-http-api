package types

import (
	"encoding/json"
)

type APIErrorResponse struct {
	Errors []string `json:"errors,omitempty"`
}

func IrrecoverableError() APIErrorResponse {
	return APIErrorResponse{Errors: []string{"irrecoverable error"}}
}

func MarshalErrors(errlist []string) []byte {
	errs := APIErrorResponse{Errors: errlist}
	res, err := json.Marshal(errs)
	if err != nil {
		res, _ = json.Marshal(IrrecoverableError())
	}
	return res
}

func MarshalError(err error) []byte {
	return MarshalErrors([]string{err.Error()})
}
