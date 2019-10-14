package v1

import (
	"net/http"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

type CreateMessageHandler struct {
	Callback func(*types.Message) (*types.Message, error)
}

func (cmh *CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newMsg := types.Message{}
	if err := handlers.RequireParams(r, "to", "value", "gasPrice", "gasLimit", "method", "parameters"); err != nil {
		handlers.Respond(w, newMsg, err)
		return
	}
	if err := newMsg.BindRequest(r); err != nil {
		handlers.Respond(w, newMsg, err)
		return
	}

	executedMsg, err := cmh.Callback(&newMsg)
	handlers.Respond(w, executedMsg, err)
}
