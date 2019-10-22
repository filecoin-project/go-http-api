package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
)

func TestSignedMessageHandler_ServeHTTP(t *testing.T) {
	uri := "http://localhost/chain/signed-messages"
	signedMessage := types.SignedMessage{
		Nonce:      10,
		From:       fixtures.TestAddress0,
		To:         fixtures.TestAddress1,
		Value:      big.NewInt(12),
		GasPrice:   big.NewInt(1),
		GasLimit:   uint64(300),
		Method:     "updatePeerID",
		Parameters: []string{},
		Signature:  "abcd1234",
	}
	t.Run("message is sent and returned with extra fields filled out if callback succeeds", func(t *testing.T) {

		h := &v1.SignedMessageHandler{Callback: happyPathSMCallback}
		jsonBody, err := json.Marshal(signedMessage)
		require.NoError(t, err)
		rr := test.PostTestRequest(uri, strings.NewReader(string(jsonBody[:])), h)
		require.NotNil(t, rr)
		assert.Equal(t, http.StatusOK, rr.Code)

		signedMessage.Kind = "signedMessage"
		signedMessage.ID = "MESSAGEID"

		var actualMsg types.SignedMessage
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actualMsg))
		assert.Equal(t, signedMessage, actualMsg)
	})

	t.Run("missing params returns an error  code 400", func(t *testing.T) {
		msg := types.SignedMessage{
			To:        fixtures.TestAddress0,
			Signature: "abcd123",
		}

		msgBody, err := json.Marshal(msg)
		require.NoError(t, err)
		h := &v1.SignedMessageHandler{Callback: happyPathSMCallback}
		rr := test.PostTestRequest(uri, strings.NewReader(string(msgBody[:])), h)

		expBody := types.MarshalErrorStrings("missing parameters: From,Nonce,Value,GasPrice,GasLimit,Method")
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, expBody, rr.Body.Bytes())
	})

	t.Run("when callback returns error, it is returned with code 500", func(t *testing.T) {
		msgBody, err := json.Marshal(signedMessage)
		require.NoError(t, err)
		h := &v1.SignedMessageHandler{Callback: sadPathSMCallback}
		rr := test.PostTestRequest(uri, strings.NewReader(string(msgBody[:])), h)

		expBody := types.MarshalErrorStrings("boom")
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, expBody, rr.Body.Bytes())
	})
}

func happyPathSMCallback(sm *types.SignedMessage) (*types.SignedMessage, error) {
	execMsg := *sm
	execMsg.Kind = "signedMessage"
	execMsg.ID = "MESSAGEID"
	return &execMsg, nil
}

func sadPathSMCallback(_ *types.SignedMessage) (*types.SignedMessage, error) {
	return nil, errors.New("boom")
}
