package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

type MessageHandler struct {
	Callback func(string) (*types.Message, error)
}

func (mh *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msgId := chi.URLParam(r, "messageId")
	var marshaled []byte

	msg, err := mh.Callback(msgId)
	if err != nil {
		marshaled = types.MarshalErrors([]string{err.Error()})
	} else {
		msg.Kind = "message"
		marshaled, _ = json.Marshal(msg)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(marshaled[:])) // nolint: errcheck
}
