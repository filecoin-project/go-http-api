package v1

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

// PostCallbackResultFunc defines a function that performs a post to an API
// Client's provided callback URL
type PostCallbackResultFunc func(cburl *url.URL, body string)

// WaitForMessageHandlerCb defines a function that calls the API implementer's
// WaitForMesssageHandler callback
type WaitForMessageHandlerCb func(*cid.Cid, *big.Int) (*types.SignedMessage, error)

// WaitForMessageHandler is the handler for the GET
type WaitForMessageHandler struct {
	Callback WaitForMessageHandlerCb
	PostCallbackResultFunc
}

// WaitForMessageParams holds the params parsed from the incoming http.Request
type WaitForMessageParams struct {
	MsgCid      *cid.Cid `json:"msgCid"`      // message Cid
	BlockHeight *big.Int `json:"blockHeight"` // block height at which to stop waiting for the message
	CbURL       *url.URL `json:"cbUrl"`       // callback URL for posting the result
}

// NewWaitForMessageHandler creates a new WaitForMessageHandler with the provided
// Filecoin Callback func and PostCallbackResultFunc
func NewWaitForMessageHandler(cb WaitForMessageHandlerCb, pcbrf PostCallbackResultFunc) *WaitForMessageHandler {
	if pcbrf == nil {
		pcbrf = PostToCallbackURL
	}
	return &WaitForMessageHandler{
		Callback:               cb,
		PostCallbackResultFunc: pcbrf,
	}
}

// ServeHTTP handles an HTTP request.
func (wfmh *WaitForMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqParams, err := wfmh.getParams(r)
	if err != nil {
		handlers.RespondBadRequest(w, err)
	}

	var msg *types.SignedMessage
	go func() {
		msg, err = wfmh.Callback(reqParams.MsgCid, reqParams.BlockHeight)
		if err != nil {
			wfmh.PostCallbackResultFunc(reqParams.CbURL, string(types.MarshalError(err)[:]))
			return
		}
		body, err := msg.MarshalJSON()
		if err != nil {
			log.Error(err)
			return
		}
		wfmh.PostCallbackResultFunc(reqParams.CbURL, string(body[:]))
	}()
	handlers.Respond(w, msg, err)
}

// getParams parses params from the http requests
func (wfmh *WaitForMessageHandler) getParams(r *http.Request) (WaitForMessageParams, error) {
	var body []byte
	_, err := r.Body.Read(body)
	if err != nil {
		return WaitForMessageParams{}, err
	}

	var rp WaitForMessageParams
	if err = json.Unmarshal(body, &rp); err != nil {
		return rp, err
	}
	return rp, nil
}

// PostToCallbackURL posts a message with the given body as JSON to the given URL.
func PostToCallbackURL(cburl *url.URL, body string) {
	resp, err := http.Post(cburl.String(), "application/json", strings.NewReader(body))
	if err != nil {
		log.Error(err)
	}
	log.Debugf("callback url response: %s %s", resp.StatusCode, resp.Body)
}
