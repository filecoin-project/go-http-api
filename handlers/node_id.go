package handlers

import (
	"fmt"
	"net/http"
)


type NodeID struct {
	Callback func()(string, error)
}

func (nid *NodeID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := nid.Callback()

	if err != nil {
		// add to JSON error struct, return that
	}

	fmt.Fprint(w, result)
}
