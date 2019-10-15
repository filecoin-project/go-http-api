package v1_test

import (
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/types"
)

func TestCreateMessageHandler_ServeHTTP(t *testing.T) {
	path := "chain/messages"
	t.Run("Message is 'created' and returned if callback succeeds", func(t *testing.T) {
		expectedMsg := types.Message{
			To:         "someaddr",
			Value:      big.NewInt(314),
			GasPrice:   big.NewInt(1),
			GasLimit:   uint64(300),
			Method:     "updatePeerID",
			Parameters: []string{"QmSTGFu1zZwrAvS8m9cWiZuuZ5pR33m85JJBuKnVPzX3u5"},
		}

		cb := func(message *types.Message) (*types.Message, error) {
			expectedMsg = types.Message{
				ID:         "sll3525ieiaghaQOEI582slkd0LKDFIeoiwRDeus",
				Nonce:      878,
				From:       "t27syykyps4puabw5fol3kn4n7flp44dz772hk3sq",
				To:         message.To,
				Value:      message.Value,
				GasPrice:   message.GasPrice,
				GasLimit:   message.GasLimit,
				Method:     message.Method,
				Parameters: message.Parameters,
				Signature:  "STcLQ6ULcLreAhwCNtsd4GICPq9AN2JGJWa315zli4NqqphgOtxK4I",
			}
			return &expectedMsg, nil
		}

		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)
		req := RequireHTTPPost(t, path, strings.NewReader(string(jsonBody[:])))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := v1.CreateMessageHandler{Callback: cb}
		handler.ServeHTTP(rr, req)

		var executedMsg types.Message
		expectedMsg.Kind = "message"

		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &executedMsg))
		assert.Equal(t, expectedMsg, executedMsg)
	})

	t.Run("if body (message) is not provided, responds with error", func(t *testing.T) {
		req := RequireHTTPPost(t, path, strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")

		req.PostForm = url.Values{}
		rr := httptest.NewRecorder()
		cb := func(*types.Message) (*types.Message, error) {
			return &types.Message{}, errors.New("should not happen")
		}
		handler := v1.CreateMessageHandler{Callback: cb}
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status == http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		// Check the response body is what we expect.
		expected := types.APIErrorResponse{Errors: []string{
			"missing message parameters",
		}}
		var actual types.APIErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))
		assert.Equal(t, expected, actual)
	})

	t.Run("if required parameters are not provided, responds with error", func(t *testing.T) {
		expectedMsg := types.Message{To: "someaddr"}
		jsonBody, err := json.Marshal(expectedMsg)
		require.NoError(t, err)
		req := RequireHTTPPost(t, path, strings.NewReader(string(jsonBody[:])))
		req.Header.Set("Content-Type", "application/json")

		req.PostForm = url.Values{}
		rr := httptest.NewRecorder()
		cb := func(*types.Message) (*types.Message, error) {
			return &types.Message{}, errors.New("should not happen")
		}
		handler := v1.CreateMessageHandler{Callback: cb}
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status == http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		// Check the response body is what we expect.
		expected := types.APIErrorResponse{Errors: []string{
			"missing parameters: GasLimit,Method",
		}}
		var actual types.APIErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))
		assert.Equal(t, expected, actual)

	})

}

func RequireHTTPPost(t *testing.T, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest("POST", "http://localhost/"+path, body)
	require.NotNil(t, req)
	return req
}
