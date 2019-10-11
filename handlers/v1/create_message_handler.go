package v1

import (
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type CreateMessageHandler struct {
	Callback func(*types.Message) (*types.Message, error)
}

func (cmh *CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := &types.MessageRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err)) //nolint: errcheck
		return
	}
	newMsg := data.Message
	executedMsg, err := cmh.Callback(newMsg)
	Respond(w, executedMsg, err)
}

func (cmh *CreateMessageHandler) ServeHTTPOld(w http.ResponseWriter, r *http.Request) {
	var value, gasPrice big.Int
	var valStr, gasPriceStr string
	var gasLimit uint64

	valStr = chi.URLParam(r, "value")
	value.SetString(valStr, 10)

	gasPriceStr = chi.URLParam(r, "gasPrice")
	gasPrice.SetString(gasPriceStr, 10)

	gasLimit, err := strconv.ParseUint(chi.URLParam(r, "gasLimit"), 10, 64)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgParams := strings.Split(chi.URLParam(r, "parameters"), ",")

	newMsg := types.Message{
		To:         chi.URLParam(r, "to"),
		Value:      &value,
		GasPrice:   &gasPrice,
		GasLimit:   gasLimit,
		Method:     chi.URLParam(r, "method"),
		Parameters: msgParams,
	}

	executedMsg, err := cmh.Callback(&newMsg)
	Respond(w, executedMsg, err)
}
