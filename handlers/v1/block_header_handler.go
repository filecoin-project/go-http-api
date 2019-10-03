package v1

import (
	"encoding/json"
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/carbonfive/go-filecoin-rest-api/types/api_errors"
	"github.com/gorilla/mux"
	"net/http"
)

// BlockHeaderHandler is a handler for the blockheader endpoint
type BlockHeaderHandler struct {
	Callback func(blockId string) (*types.BlockHeader, error)
}

// ServeHTTP handles an HTTP request.
func (b *BlockHeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var marshaled []byte

	result, err := b.Callback(vars["blockId"])
	if err != nil {
		marshaled = api_errors.MarshalErrors([]string{err.Error()})
	} else {
		marshaled, _ = json.Marshal(result)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(marshaled[:])) // nolint: errcheck
}
