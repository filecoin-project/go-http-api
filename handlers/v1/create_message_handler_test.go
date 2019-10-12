package v1_test

import (
	"encoding/json"
	"errors"
	v1 "github.com/carbonfive/go-filecoin-rest-api/handlers/v1"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateMessageHandler_ServeHTTP(t *testing.T) {
	path := "chain/messages"
	t.Run("Message is 'created' and returned if callback succeeds", func(t *testing.T) {
		to := "someaddr"
		value := "314"
		gasPrice := "1"
		gasLimit := "300"
		method := "updatePeerID"
		msgParams := []string{"QmSTGFu1zZwrAvS8m9cWiZuuZ5pR33m85JJBuKnVPzX3u5"}
		postParams := url.Values{}
		postParams.Add("to", to)
		postParams.Add("value", value)
		postParams.Add("gasPrice", gasPrice)
		postParams.Add("gasLimit", gasLimit)
		postParams.Add("method", method)
		postParams.Add("parameters", strings.Join(msgParams, ","))

		var expectedMsg types.Message

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

		req := RequireHTTPRequest(t, path)
		req.PostForm = postParams
		rr := httptest.NewRecorder()

		handler := v1.CreateMessageHandler{Callback: cb}
		handler.ServeHTTP(rr, req)

		var executedMsg types.Message
		expectedMsg.Kind = "message"

		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &executedMsg))
		assert.Equal(t, expectedMsg, executedMsg)
	})

	t.Run("if required params are not provided, responds with error", func(t *testing.T) {
		req := RequireHTTPRequest(t, path)

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
			"to, value, gasPrice, gasLimit, method, parameters,  required parameters are missing",
		}}
		var actual types.APIErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))
		assert.Equal(t, expected, actual)
	})
}

func RequireHTTPRequest(t *testing.T, path string) *http.Request {
	req, err := http.NewRequest("POST", path, nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	return req
}
