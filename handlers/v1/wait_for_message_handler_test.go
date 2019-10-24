package v1_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"

	. "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestNewWaitForMessageHandler(t *testing.T) {
	cb := func(*cid.Cid, *big.Int) (*types.SignedMessage, error) {
		return &types.SignedMessage{}, nil
	}

	t.Run("sets provided callback func", func(t *testing.T) {
		wfmh := NewWaitForMessageHandler(cb)
		AssertEqualFuncs(t, cb, wfmh.Callback)
	})

}

func TestWaitForMessageHandler_ServeHTTP(t *testing.T) {
	cid1 := test.RequireTestCID(t, []byte("cid1"))
	cburl := "http://bigmoney-nowhammies.com/message-complete"
	bh := big.NewInt(8)

	uri := "http://localhost:5000/chain/messages/abcd1234/wait"
	h := WaitForMessageHandler{
		Callback: func(_ *cid.Cid, _ *big.Int) (*types.SignedMessage, error) {
			return &types.SignedMessage{}, nil
		},
	}

	t.Run("Returns SignedMessage struct", func(t *testing.T) {
		params := &[]test.Param{
			{Key: "msgCid", Value: cid1.String()},
			{Key: "blockHeight", Value: bh.String()},
			{Key: "callbackURL", Value: cburl},
		}
		rr := test.GetTestRequest(uri, params, &h)
		assert.Equal(t, http.StatusOK, rr.Code)

		var expMsg types.SignedMessage
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &expMsg))
		assert.Equal(t, "signedMessage", expMsg.Kind)
	})

	t.Run("If msgCid fails to decode, returns error", func(t *testing.T) {
		params := &[]test.Param{
			{Key: "msgCid", Value: "not valid"},
			{Key: "blockHeight", Value: bh.String()},
		}
		rr := test.GetTestRequest(uri, params, &h)
		expErr := `{"errors":["msgCid 'not valid': selected encoding not supported"]}`

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())
	})

	t.Run("If blockHeight fails to unmarshal, returns error", func(t *testing.T) {
		params := &[]test.Param{
			{Key: "msgCid", Value: cid1.String()},
			{Key: "blockHeight", Value: "not valid"},
		}
		expErr := `{"errors":["blockHeight 'not valid': failed to parse"]}`
		rr := test.GetTestRequest(uri, params, &h)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())
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
