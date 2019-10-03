package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/types/api_errors"
	"net/http"
)

type ActorsHandler struct {
	Callback func() (json []byte, err error)
}

// ServeHTTP handles an HTTP request.
func (a *ActorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := a.Callback()

	w.WriteHeader(http.StatusOK)
	if err != nil {
		result = api_errors.MarshalErrors([]string{err.Error()})
	}

	fmt.Fprint(w, string(result[:])) // nolint: errcheck
}
