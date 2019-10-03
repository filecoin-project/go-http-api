package v1

import (
	"fmt"
	"net/http"

	"github.com/carbonfive/go-filecoin-rest-api/types/api_errors"
)

// NodeHandler is the handler for the control/node endpoint
type NodeHandler struct {
	Callback func() (json []byte, err error)
}

func (nid *NodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := nid.Callback()

	if err != nil {
		result = api_errors.MarshalErrors([]string{err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(result[:])) // nolint: errcheck
}
