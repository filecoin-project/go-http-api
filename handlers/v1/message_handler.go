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
	var marshaled []byte
	msgId := chi.URLParam(r, "messageId")

	msg, err := mh.Callback(msgId)
	if err != nil {
		marshaled = types.MarshalErrors([]string{err.Error()})
	} else {
		msg.Kind = "message"
		if marshaled, err = json.Marshal(msg) ; err != nil {
			log.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _,err = fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
