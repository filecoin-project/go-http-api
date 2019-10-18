package v1_test

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
	"github.com/stretchr/testify/require"
)

func TestSignedMessageHandler_ServeHTTP(t *testing.T) {
	uri := "http://localhost/chain/signed-messages"
	t.Run("message is sent and returned with extra fields filled out if callback succeeds", func(t *testing.T) {
		expectedMsg := types.SignedMessage{
			Message: types.Message{
				Kind:       "message",
				Nonce:      10,
				From:       fixtures.TestAddress0,
				To:         fixtures.TestAddress1,
				Value:      big.NewInt(12),
				GasPrice:   big.NewInt(1),
				GasLimit:   uint64(300),
				Method:     "updatePeerID",
				Parameters: []string{},
			},
			Signature: "abcd1234",
		}
		h := &v1.SignedMessageHandler{Callback: happyPathSMCallback}
		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)
		rr := test.PostTestRequest(uri, strings.NewReader(string(jsonBody[:])), h)
		require.NotNil(t, rr)
	})

	t.Run("missing params returns an error", func(t *testing.T) {

	})
}

func happyPathSMCallback(sm *types.SignedMessage) (*types.SignedMessage, error) {
	sm.ID = "SLDFKJSLDFKJSDLFKJSDFLMESSAGEID"
	return sm, nil
}

//func sadPathSMCallback(message *types.SignedMessage) (*types.SignedMessage, error) {
//	return nil, errors.New("boom")
//}
