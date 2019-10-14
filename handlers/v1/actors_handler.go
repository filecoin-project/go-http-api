package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-http-api/types"
)

// ActorsHandler is the handler for the actors endpoint
type ActorsHandler struct {
	Callback func() ([]*types.Actor, error)
}

// ServeHTTP handles an HTTP request.
func (a *ActorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var marshaled []byte
	result, err := a.Callback()

	if err != nil {
		marshaled = types.MarshalErrors([]string{err.Error()})
	} else {
		for _, el := range result {
			el.Kind = "actor"
		}
		if marshaled, err = json.Marshal(result); err != nil {
			log.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err = fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
