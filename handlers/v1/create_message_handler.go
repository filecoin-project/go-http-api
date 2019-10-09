package v1

import (
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/go-chi/chi"
	"go/ast"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type CreateMessageHandler struct {
	Callback func(*types.Message)(*types.Message, error)
}

func (cmh *CreateMessageHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var value, gasPrice big.Int
	valStr := chi.URLParam(r, "value")
	value.SetString(valStr, 10)

	gasPriceStr := chi.URLParam(r, "gasPrice")
	gasPrice.SetString(gasPriceStr, 10)

	gasLimit, err :=  strconv.ParseUint(chi.URLParam(r, "gasLimit"), 10, 64)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parameters := strings.Split(chi.URLParam(r, "parameters"),",")

	newMsg := types.Message{
		To:         chi.URLParam(r, "to"),
		Value:      &value,
		GasPrice:   &gasPrice,
		GasLimit:   gasLimit,
		Method:     chi.URLParam(r,"method"),
		Parameters: parameters,
	}

	executedMsg, err := cmh.Callback(&newMsg)
	Respond(w, executedMsg, err)
}
