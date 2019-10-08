package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

// NodeHandler is the handler for the control/node endpoint
type NodeHandler struct {
	Callback func() (*types.Node, error)
}

func (nid *NodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var marshaled []byte
	node, err := nid.Callback()

	if err != nil {
		marshaled = types.MarshalErrors([]string{err.Error()})
	} else {
		node.Kind= "node"
		if marshaled, err = json.Marshal(node); err != nil {
			log.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _,err = fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
