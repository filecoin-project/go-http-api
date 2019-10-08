package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

// BlockHandler is a handler for the blockheader endpoint
type BlockHandler struct {
	Callback func(blockId string) (*types.Block, error)
}

// ServeHTTP handles an HTTP request.
func (bhh *BlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var marshaled []byte
	blockId := chi.URLParam(r, "actorId")

	block, err := bhh.Callback(blockId)
	if err != nil {
		marshaled = types.MarshalErrors([]string{err.Error()})
	} else {
		block.Kind = "block"
		block.Header.Kind = "blockHeader"
		if marshaled, err = json.Marshal(block); err != nil {
			log.Error(err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	if _,err = fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
