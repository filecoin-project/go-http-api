package v1

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

// MessageHandler is the handler for the /chain/executed-messages/{executedMessageId} endpoint
type MessageHandler struct {
	Callback func(string) (*types.SignedMessage, error)
}

// ServeHTTP handles an HTTP request.
func (mh *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msgID := chi.URLParam(r, "messageId")

	msg, err := mh.Callback(msgID)
	handlers.Respond(w, msg, err)
}
