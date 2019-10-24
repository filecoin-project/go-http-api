package v1_test

import (
	"fmt"
	"github.com/filecoin-project/go-http-api/test"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"

	. "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/types"
)

func TestNewWaitForMessageHandler(t *testing.T) {
	pcb := func(_ *url.URL, _ string) { return }
	cb := func(*cid.Cid, *big.Int) (*types.SignedMessage, error) {
		return &types.SignedMessage{}, nil
	}

	t.Run("sets provided callback & post callback funcs", func(t *testing.T) {
		wfmh := NewWaitForMessageHandler(cb, pcb)
		AssertEqualFuncs(t, cb, wfmh.Callback)
		AssertEqualFuncs(t, pcb, wfmh.PostCallbackResultFunc)
	})

	t.Run("sets default post callback func if param is nil", func(t *testing.T) {
		wfmh := NewWaitForMessageHandler(cb, nil)
		AssertEqualFuncs(t, PostToCallbackURL, wfmh.PostCallbackResultFunc)
	})
}

type testCbs struct {
	postCbWasCalled bool
}

func (t testCbs) PostCallbackTest(cburl *url.URL, body string) {
	t.postCbWasCalled = true
	return
}

func TestWaitForMessageHandler_ServeHTTP(t *testing.T) {
	cid1 := test.RequireTestCID(t, []byte("cid1"))
	cburl := "http://bigmoney-nowhammies.com/message-complete"
	bh := big.NewInt(8)

	tcbs := testCbs{}

	h := WaitForMessageHandler{
		PostCallbackResultFunc: tcbs.PostCallbackTest,
		Callback: func(_ *cid.Cid, _ *big.Int) (*types.SignedMessage, error) {
			return &types.SignedMessage{}, nil
		},
	}
	rr := httptest.NewRecorder()

	t.Run("If msgCid fails to decode, returns error", func(t *testing.T) {
		msgParams := fmt.Sprintf(`{"msgCid":"%s","blockHeight":%s,"cbUrl":"%s"}`, "not valid", bh.String(), cburl)
		req := httptest.NewRequest("GET", "http://localhost:5000/doesntmatter", strings.NewReader(msgParams))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rr, req)
	})

	t.Run("If blockHeight fails to unmarshal, returns error", func(t *testing.T) {
		msgParams := fmt.Sprintf(`{"msgCid":"%s","blockHeight":%s,"cbUrl":"%s"}`, cid1, "not valid", cburl)
		req := httptest.NewRequest("GET", "http://localhost:5000/doesntmatter", strings.NewReader(msgParams))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("if callbackURL fails to parse, returns error", func(t *testing.T) {
		msgParams := fmt.Sprintf(`{"msgCid":"%s","blockHeight":%s,"cbUrl":"%s"}`, cid1, bh.String(), "not valid")
		req := httptest.NewRequest("GET", "http://localhost:5000/doesntmatter", strings.NewReader(msgParams))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestPostToCallbackURL(t *testing.T) {

}

func AssertEqualFuncs(t *testing.T, fn1, fn2 interface{}) {
	assert.Equal(t, FuncPtrAsString(fn1), FuncPtrAsString(fn2))
}

func FuncPtrAsString(fn interface{}) string {
	res := fmt.Sprintf("%v", fn)
	return res
}
