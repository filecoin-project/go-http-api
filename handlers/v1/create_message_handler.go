package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

type CreateMessageHandler struct {
	Callback func(*types.Message) (*types.Message, error)
}

func (cmh *CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newMsg types.Message

	if err := render.Bind(r, &newMsg); err != nil {
		if err.Error() == "EOF" {
			err = errors.New("missing message parameters")
		}
		handlers.RespondBadRequest(w, err)
		return
	}

	executedMsg, err := cmh.Callback(&newMsg)
	handlers.Respond(w, executedMsg, err)
}
