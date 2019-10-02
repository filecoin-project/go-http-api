package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers/api_errors"
	"github.com/gorilla/mux"
	"net/http"
)

type Block struct {
	Callback func(blockId string) (json []byte, err error)
}

// ServeHTTP handles an HTTP request.
func (b *Block) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result, err := b.Callback(vars["blockId"])

	if err != nil {
		result = api_errors.MarshalErrors([]string{err.Error()})
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, result) // nolint: errcheck
}
