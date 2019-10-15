package v1

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

type MessageHandler struct {
	Callback func(string) (*types.Message, error)
}

func (mh *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "messageId")

	msg, err := mh.Callback(msgId)
	handlers.Respond(w, msg, err)
}
