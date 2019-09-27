package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Actor struct {
	Callback func(actorId string) (json []byte, err error)
}

func (a *Actor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result, err := a.Callback(vars["actorId"])

	w.WriteHeader(http.StatusOK)
	if err != nil {
		// add to JSON error struct, return that
		fmt.Fprint(w, handlers.MarshalError(err)) // nolint: errcheck
		return
	}

	fmt.Fprint(w, string(result[:])) // nolint: errcheck
}
