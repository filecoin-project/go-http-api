package v1

import (
	"github.com/filecoin-project/go-http-api/handlers"
	"net/http"

	"github.com/filecoin-project/go-http-api/types"
)

// NodeHandler is the handler for the control/node endpoint
type NodeHandler struct {
	Callback func() (*types.Node, error)
}

func (nid *NodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	node, err := nid.Callback()
	handlers.Respond(w, node, err)
}
