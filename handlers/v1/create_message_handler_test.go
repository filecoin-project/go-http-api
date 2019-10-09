package v1_test

import (
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateMessageHandler_ServeHTTP(t *testing.T) {
	to := "someaddr"
	value := "3.14"
	gasPrice := "1"
	gasLimit := "300"
	method := "updatePeerID"
	parameters := []string{"QmSTGFu1zZwrAvS8m9cWiZuuZ5pR33m85JJBuKnVPzX3u5"}

	cmhcb := func(message *types.Message) (*types.Message, error) {
		execMsg := types.Message{
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
		return &execMsg, nil
	}

	s := test.CreateTestServer(t, &server.V1Callbacks{CreateMessage:cmhcb}, false).Run()
	defer func() {
		assert.NoError(t, s.Shutdown())
	}()

	body := test.RequireGetResponseBody(t, s.Config().Port, "")
}
