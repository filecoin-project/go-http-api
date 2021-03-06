package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestMessageHandler_ServeHTTP(t *testing.T) {
	t.Run("GetMessageByID accepts and uses param messageId", func(t *testing.T) {
		expected := types.SignedMessage{
			ID:       "someid",
			Nonce:    98348,
			From:     "abcd1234",
			To:       "1234abcd",
			Value:    big.NewInt(8383),
			GasPrice: big.NewInt(3432),
			GasLimit: 10,
			Method:   "createMiner",
		}
		testcb := func(msgId string) (*types.SignedMessage, error) {
			return &expected, nil
		}

		s := test.CreateTestServer(t, &v1.Callbacks{GetMessageByID: testcb}, false).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "chain/executed-messages/someid")
		expected.Kind = "signedMessage"
		var actual types.SignedMessage
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, expected, actual)
	})

	t.Run("GetMessageByID passes on errors returned by Callback", func(t *testing.T) {
		expected := errors.New("boom")
		testcb := func(msgId string) (*types.SignedMessage, error) {
			return nil, expected
		}

		h := &v1.MessageHandler{Callback: testcb}
		rr := test.GetTestRequest("http://localhost:3000/chain/executed-messages/1234", nil, h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		exp := types.MarshalError(expected)
		assert.Equal(t, exp, rr.Body.Bytes())
	})
}
