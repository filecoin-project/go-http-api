package v1

import (
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"net/http"
)

type CreateMessageHandler struct {
	Callback func(*types.Message) (*types.Message, error)
}

func (cmh *CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newMsg := types.Message{}
	if err := newMsg.BindRequest(r); err != nil {
		Respond(w, newMsg, err)
	}

	executedMsg, err := cmh.Callback(&newMsg)
	Respond(w, executedMsg, err)
}
