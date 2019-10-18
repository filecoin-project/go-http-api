package types

import (
	"encoding/json"
)

// APIErrorResponse is a struct for returning errors in the API
type APIErrorResponse struct {
	Errors []string `json:"errors,omitempty"`
}

// IrrecoverableError is an error message for when even a simple statement fails
func IrrecoverableError() APIErrorResponse {
	return APIErrorResponse{Errors: []string{"irrecoverable error"}}
}

// MarshalErrorStrings takes a list of strings, and returns a marshaled APIErrorResponse
func MarshalErrorStrings(errlist ...string) []byte {
	errs := APIErrorResponse{Errors: errlist}
	res, err := json.Marshal(errs)
	if err != nil {
		res, _ = json.Marshal(IrrecoverableError())
	}
	return res
}

// MarshalError takes an error and returns a marshaled APIErrorResponse
func MarshalError(err error) []byte {
	return MarshalErrorStrings(err.Error())
}
