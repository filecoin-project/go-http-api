package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers/api_errors"
	"net/http"
)

type Node struct {
	Callback func() (json []byte, err error)
}

func (nid *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := nid.Callback()

	if err != nil {
		result = api_errors.MarshalErrors([]string{err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(result[:])) // nolint: errcheck
}
