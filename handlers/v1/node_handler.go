package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/types/api_errors"
	"net/http"
)

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
