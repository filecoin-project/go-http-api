package types

import (
	"encoding/json"
)

// APIErrorResponse is a struct for returning errors in the API
type APIErrorResponse struct {
	Errors []string `json:"errors,omitempty"`
}

// MarshalErrorStrings takes a list of strings, and returns a marshaled APIErrorResponse
func MarshalErrorStrings(errlist ...string) []byte {
	errs := APIErrorResponse{Errors: errlist}
	res, _ := json.Marshal(errs) // nolint: errcheck
	return res
}

// MarshalError takes an error and returns a marshaled APIErrorResponse
func MarshalError(err error) []byte {
	return MarshalErrorStrings(err.Error())
}
