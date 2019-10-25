package v1_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-http-api/test/fixtures"
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
		test.AssertEqualFuncs(t, cb, wfmh.Callback)
	})

}

func TestWaitForMessageHandler_ServeHTTP(t *testing.T) {
	cid1 := test.RequireTestCID(t, []byte("cid1"))
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

	t.Run("If callback fails, returns error", func(t *testing.T) {
		badHandler := WaitForMessageHandler{Callback: func(_ *cid.Cid, _ *big.Int) (message *types.SignedMessage, e error) {
			return nil, errors.New("boom")
		}}
		params := &[]test.Param{
			{Key: "msgCid", Value: cid1.String()},
			{Key: "blockHeight", Value: bh.String()},
		}
		rr := test.GetTestRequest(uri, params, &badHandler)
		expErr := `{"errors":["boom"]}`
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())

	})

	t.Run("endpoint can be called", func(t *testing.T) {
		cid1 := test.RequireTestCID(t, []byte("cid1"))
		bh := big.NewInt(8)
		expMsg := types.SignedMessage{
			ID:        cid1.String(),
			Nonce:     10,
			From:      fixtures.TestAddress0,
			To:        fixtures.TestAddress1,
			Value:     big.NewInt(0),
			GasPrice:  big.NewInt(0),
			GasLimit:  0,
			Method:    "updatePeerID",
			Signature: "somesig",
		}
		cb := func(_ *cid.Cid, _ *big.Int) (*types.SignedMessage, error) {
			return &expMsg, nil
		}
		cbs := Callbacks{WaitForMessage: cb}
		path := fmt.Sprintf("chain/messages/%s/wait?blockHeight=%s", cid1.String(), bh.String())
		test.AssertServerResponse(t, &cbs, false, path, "foo")
	})
}
