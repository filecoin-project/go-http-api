package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-http-api/types"
)

// Respond is a standardized response for handlers
func Respond(w http.ResponseWriter, result interface{}, cberr error) {
	var marshaled []byte
	var err error

	if cberr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		marshaled = types.MarshalError(cberr)
	} else {
		if marshaled, err = json.Marshal(result); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	if _, err := fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}

// RespondBadRequest is a standardized response for API errors where
// the request was malformed
func RespondBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	marshaled := types.MarshalError(err)
	if _, err := fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
