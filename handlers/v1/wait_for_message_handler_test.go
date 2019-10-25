package v1_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
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
	// when using httptest the URL doesn't matter; we're just testing
	// the handler
	uri := "http://localhost:5000/doesntmatter"

	h := WaitForMessageHandler{
		Callback: func(_ *cid.Cid, _ *big.Int) (*types.SignedMessage, error) {
			return &types.SignedMessage{}, nil
		},
	}

	t.Run("Returns SignedMessage struct", func(t *testing.T) {
		params := url.Values{}
		params.Set("blockHeight", bh.String())
		// but we must pass everything so chi url params are set for httptest
		params.Set("messageCid", cid1.String())

		rr := test.GetTestRequest(uri, params, &h)
		assert.Equal(t, http.StatusOK, rr.Code)

		var expMsg types.SignedMessage
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &expMsg))
		assert.Equal(t, "signedMessage", expMsg.Kind)
	})

	t.Run("If messageCid fails to decode, returns error", func(t *testing.T) {
		params := url.Values{}
		params.Set("blockHeight", bh.String())
		// have to do this so chi url params are set for httptest
		params.Set("messageCid", "notvalid")

		rr := test.GetTestRequest(uri, params, &h)
		expErr := `{"errors":["messageCid 'notvalid': selected encoding not supported"]}`

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())
	})

	t.Run("If blockHeight fails to unmarshal, returns error", func(t *testing.T) {
		params := url.Values{}
		params.Set("blockHeight", "not valid")
		// have to do this so chi url params are set for httptest
		params.Set("messageCid", cid1.String())

		expErr := `{"errors":["blockHeight 'not valid': failed to parse"]}`
		rr := test.GetTestRequest(uri, params, &h)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())
	})

	t.Run("If callback fails, returns error", func(t *testing.T) {
		badHandler := WaitForMessageHandler{Callback: func(_ *cid.Cid, _ *big.Int) (message *types.SignedMessage, e error) {
			return nil, errors.New("boom")
		}}
		params := url.Values{}
		params.Set("blockHeight", bh.String())
		params.Set("messageCid", cid1.String())

		rr := test.GetTestRequest(uri, params, &badHandler)
		expErr := `{"errors":["boom"]}`
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, expErr, rr.Body.String())
	})
}

// Test round trip so we know routing is working
func TestIntegration_ServeHTTP(t *testing.T) {
	cid1 := test.RequireTestCID(t, []byte("cid1"))
	bh := big.NewInt(8)

	t.Run("endpoint can be called", func(t *testing.T) {
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

		s := test.CreateTestServer(t, &cbs, false)
		s.Run()
		defer s.Shutdown() // nolint: errcheck

		path := fmt.Sprintf("chain/messages/%s/wait", cid1.String())
		requri := fmt.Sprintf("http://localhost:%d/api/filecoin/v1/%s?", s.Config().Port, path)

		params := url.Values{}
		params.Set("blockHeight", bh.String())

		resp, err := http.Get(requri + params.Encode())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		expMsg.Kind = "signedMessage"
		test.AssertMarshaledEquals(t, &expMsg, string(body[:]))
	})
}
