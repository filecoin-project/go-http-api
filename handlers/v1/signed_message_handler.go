package v1

import (
	"errors"
	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
	"github.com/go-chi/render"
	"net/http"
)

type SignedMessageHandler struct {
	Callback func(*types.SignedMessage) (*types.SignedMessage, error)
}

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
