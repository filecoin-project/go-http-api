package v1

import (
	"fmt"
	"net/http"
)

type Node struct {
	Callback func() (json []byte, err error)
}

func (nid *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := nid.Callback()

	if err != nil {
		// add to JSON error struct, return that
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(result[:])) // nolint: errcheck
}
