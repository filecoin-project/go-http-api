package v1_test

import (
	"encoding/json"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
	"testing"
)

func TestCreateMessageHandler_ServeHTTP(t *testing.T) {
	path := "chain/messages"

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

	cmhcb := func(message *types.Message) (*types.Message, error) {
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

	s := test.CreateTestServer(t, &server.V1Callbacks{CreateMessage: cmhcb}, false).Run()
	defer func() {
		assert.NoError(t, s.Shutdown())
	}()

	body := test.RequirePostFormResponseBody(t, s.Config().Port, path, postParams)

	var executedMsg types.Message

	require.NoError(t, json.Unmarshal(body, &executedMsg))
	assert.Equal(t, "message", executedMsg.Kind)
	assert.Equal(t, expectedMsg.Value, executedMsg.Value)
	assert.Equal(t, expectedMsg.ID, executedMsg.ID)
	assert.Equal(t, expectedMsg.Nonce, executedMsg.Nonce)
	assert.Equal(t, expectedMsg.GasPrice, executedMsg.GasPrice)
	assert.Equal(t, expectedMsg.GasLimit, executedMsg.GasLimit)
	assert.Equal(t, expectedMsg.Method, executedMsg.Method)
	assert.Equal(t, expectedMsg.Parameters, executedMsg.Parameters)
	assert.Equal(t, expectedMsg.Signature, executedMsg.Signature)
}
