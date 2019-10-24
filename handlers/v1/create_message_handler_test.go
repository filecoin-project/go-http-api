package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestCreateMessageHandler_ServeHTTP(t *testing.T) {
	uri := "http://localhost/chain/messages"
	expectedMsg := types.Message{
		To:         "someaddr",
		Value:      big.NewInt(314),
		GasPrice:   big.NewInt(1),
		GasLimit:   uint64(300),
		Method:     "updatePeerID",
		Parameters: []string{"QmSTGFu1zZwrAvS8m9cWiZuuZ5pR33m85JJBuKnVPzX3u5"},
	}

	t.Run("Message is 'created' and returned if callback succeeds", func(t *testing.T) {
		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", uri, strings.NewReader(string(jsonBody[:])))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := v1.CreateMessageHandler{Callback: happyPathCMCallback}
		handler.ServeHTTP(rr, req)

		var executedMsg types.Message
		expectedMsg.Kind = "message"
		expectedMsg.ID = "sll3525ieiaghaQOEI582slkd0LKDFIeoiwRDeus"
		expectedMsg.From = "t27syykyps4puabw5fol3kn4n7flp44dz772hk3sq"
		expectedMsg.Nonce = 878

		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &executedMsg))
		assert.Equal(t, expectedMsg, executedMsg)
	})

	t.Run("if callback fails, responds with error", func(t *testing.T) {
		h := &v1.CreateMessageHandler{Callback: sadPathCMCallback}

		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)

		rr := test.PostTestRequest(uri, strings.NewReader(string(jsonBody[:])), h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		expected := types.MarshalErrorStrings("boom")

		assert.Equal(t, expected, rr.Body.Bytes())

	})

	t.Run("if body (message) is not provided, responds with error", func(t *testing.T) {
		h := &v1.CreateMessageHandler{Callback: happyPathCMCallback}
		rr := test.PostTestRequest(uri, strings.NewReader(""), h)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expected := types.MarshalErrorStrings("missing message parameters")
		assert.Equal(t, expected, rr.Body.Bytes())
	})

	t.Run("if required parameters are not provided, responds with error", func(t *testing.T) {
		expectedMsg := types.Message{To: "someaddr"}
		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)
		req := httptest.NewRequest("POST", uri, strings.NewReader(string(jsonBody[:])))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h := v1.CreateMessageHandler{Callback: happyPathCMCallback}
		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expected := types.APIErrorResponse{Errors: []string{
			"missing parameters: Value,GasPrice,GasLimit,Method",
		}}
		var actual types.APIErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))
		assert.Equal(t, expected, actual)
	})

}

func happyPathCMCallback(message *types.Message) (*types.Message, error) {
	msg := types.Message{
		ID:         "sll3525ieiaghaQOEI582slkd0LKDFIeoiwRDeus",
		Nonce:      878,
		From:       "t27syykyps4puabw5fol3kn4n7flp44dz772hk3sq",
		To:         message.To,
		Value:      message.Value,
		GasPrice:   message.GasPrice,
		GasLimit:   message.GasLimit,
		Method:     message.Method,
		Parameters: message.Parameters,
	}
	return &msg, nil
}

func sadPathCMCallback(message *types.Message) (*types.Message, error) {
	return &types.Message{}, errors.New("boom")

}
