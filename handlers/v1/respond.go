package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

func Respond(w http.ResponseWriter, result interface{}, cberr error) {
	var marshaled []byte
	var err error

	if cberr != nil {
		marshaled = types.MarshalErrors([]string{cberr.Error()})
	} else {
		if marshaled, err = json.Marshal(result); err != nil {
			log.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, string(marshaled[:])); err != nil {
		log.Error(err)
	}
}
