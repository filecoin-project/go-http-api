package v1_test

import (
	"encoding/json"
	"errors"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestMessageHandler_ServeHTTP(t *testing.T) {
	t.Run("GetMessageByID accepts and uses param messageId", func(t *testing.T) {
		expected := types.Message{
			ID:         "someid",
			Nonce:      98348,
			From:       "abcd1234",
			To:         "1234abcd",
			Value:      big.NewInt(8383),
			GasPrice:   big.NewInt(3432),
			GasLimit:   10,
			Method:     "createMiner",
			Parameters: nil,
			Signature:  "somesig",
		}
		testcb := func(msgId string) (*types.Message, error) {
			return &expected, nil
		}

		s := test.CreateTestServer(t, &server.V1Callbacks{GetMessageByID: testcb}, false).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "chain/executed-messages/someid")
		var actual types.Message
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, "message", actual.Kind)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Nonce, actual.Nonce)
		assert.Equal(t, expected.From, actual.From)
		assert.Equal(t, expected.To, actual.To)
		assert.Equal(t, expected.Value, actual.Value)
		assert.Equal(t, expected.GasPrice, actual.GasPrice)
		assert.Equal(t, expected.GasLimit, actual.GasLimit)
		assert.Equal(t, expected.Method, actual.Method)
		assert.Equal(t, expected.Signature, actual.Signature)
	})

	t.Run("GetMessageByID passes on errors returned by Callback", func(t *testing.T) {
		expected := errors.New("boom!")
		testcb := func(msgId string) (*types.Message, error) {
			return nil, expected
		}

		s := test.CreateTestServer(t, &server.V1Callbacks{GetMessageByID: testcb}, false).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "chain/executed-messages/someid")
		var actual types.Message
		var apiErr types.APIErrorResponse

		undefMsg := types.Message{}
		assert.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, undefMsg, actual)

		assert.NoError(t, json.Unmarshal(body, &apiErr))
		assert.Len(t, apiErr.Errors, 1)
		assert.Equal(t, "boom!", apiErr.Errors[0])
	})
}
