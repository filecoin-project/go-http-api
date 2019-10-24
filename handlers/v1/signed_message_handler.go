package v1

import (
	"errors"
	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
	"github.com/go-chi/render"
	"net/http"
)

// SignedMessageHandler is the handler for the POST /chain/messages/ endpoint
type SignedMessageHandler struct {
	Callback func(*types.SignedMessage) (*types.SignedMessage, error)
}

// ServeHTTP handles an HTTP request.
func (ssmh *SignedMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newMsg types.SignedMessage

	if err := render.Bind(r, &newMsg); err != nil {
		if err.Error() == "EOF" {
			// message body was blank
			err = errors.New("missing signed message")
		}
		handlers.RespondBadRequest(w, err)
		return
	}
	executedMsg, err := ssmh.Callback(&newMsg)
	handlers.Respond(w, executedMsg, err)
}
