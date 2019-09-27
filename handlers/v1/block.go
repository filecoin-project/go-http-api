package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Block struct {
	Callback func(blockId string) (json string, err error)
}

func (b *Block) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result, err := b.Callback(vars["blockId"])

	if err != nil {
		// add to JSON error struct, print that
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, result) // nolint: errcheck
}
