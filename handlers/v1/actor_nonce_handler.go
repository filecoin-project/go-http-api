package v1

import (
	"math/big"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/filecoin-project/go-http-api/handlers"
)

type ActorNonceHandler struct {
	Callback func(string) (*big.Int, error)
}

func (anh *ActorNonceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actorId := chi.URLParam(r, "actorId")

	actor, err := anh.Callback(actorId)
	handlers.Respond(w, actor, err)
}
