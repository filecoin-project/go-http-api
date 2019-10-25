package v1

import (
	"fmt"
	"github.com/go-chi/chi"
	"math/big"
	"net/http"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

// WaitForMessageHandlerCb defines a function that calls the API implementer's
// WaitForMesssageHandler callback
type WaitForMessageHandlerCb func(*cid.Cid, *big.Int) (*types.SignedMessage, error)

// WaitForMessageHandler is the handler for the GET /chain/messages/{messageID}/wait
type WaitForMessageHandler struct {
	Callback WaitForMessageHandlerCb
}

// WaitForMessageParams holds the params parsed from the incoming http.Request
type WaitForMessageParams struct {
	MsgCid      *cid.Cid `json:"msgCid"`      // message Cid
	BlockHeight *big.Int `json:"blockHeight"` // block height at which to stop waiting for the message
}

// NewWaitForMessageHandler creates a new WaitForMessageHandler with the provided
// Filecoin Callback func
func NewWaitForMessageHandler(cb WaitForMessageHandlerCb) *WaitForMessageHandler {
	return &WaitForMessageHandler{Callback: cb}
}

// ServeHTTP handles an HTTP request.
func (wfmh *WaitForMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqParams, err := wfmh.getParams(r)
	if err != nil {
		handlers.RespondBadRequest(w, err)
		return
	}

	msg, err := wfmh.Callback(reqParams.MsgCid, reqParams.BlockHeight)
	handlers.Respond(w, msg, err)
}

// getParams parses needed values from an http request
func (wfmh *WaitForMessageHandler) getParams(r *http.Request) (WaitForMessageParams, error) {
	mc := chi.URLParam(r, "messageCid")
	msgCid, err := cid.Decode(mc)
	if err != nil {
		return WaitForMessageParams{}, fmt.Errorf("messageCid '%s': %s", mc, err.Error())
	}

	bh := r.FormValue("blockHeight")
	blockH, ok := big.NewInt(0).SetString(bh, 10)
	if !ok {
		return WaitForMessageParams{}, fmt.Errorf("blockHeight '%s': failed to parse", bh)
	}
	return WaitForMessageParams{
		MsgCid:      &msgCid,
		BlockHeight: blockH,
	}, nil
}
