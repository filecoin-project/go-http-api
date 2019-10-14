package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-http-api/types"
)

func Respond(w http.ResponseWriter, result interface{}, cberr error) {
	var marshaled []byte
	var err error

	if cberr != nil {
		w.WriteHeader(http.StatusBadRequest)
		marshaled = types.MarshalErrors([]string{cberr.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
		if marshaled, err = json.Marshal(result); err != nil {
			log.Error(err)
			return
		}
	}

	if _, err := fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}