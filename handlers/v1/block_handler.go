package v1

import (
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
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
	blockId := chi.URLParam(r, "actorId")
	block, err := bhh.Callback(blockId)
	handlers.Respond(w, block, err)
}
