package types_test

import (
	"math/big"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-http-api/test/fixtures"
	. "github.com/filecoin-project/go-http-api/types"
)

func TestSignedMessage_Bind(t *testing.T) {
	req := httptest.NewRequest("GET", "/some/path", nil)

	t.Run("call on full struct returns no error", func(t *testing.T) {
		sm := SignedMessage{
			Nonce:      1010,
			From:       fixtures.TestAddress0,
			To:         fixtures.TestAddress1,
			Value:      big.NewInt(1),
			GasPrice:   big.NewInt(2),
			GasLimit:   3,
			Method:     "reclaim",
			Parameters: []string{"channelID"},
			Signature:  "somesignature",
		}
		assert.NoError(t, sm.Bind(req))
	})
	t.Run("call on empty struct returns missing fields in error message", func(t *testing.T) {
		expErr := "missing parameters: To,From,Nonce,Value,GasPrice,GasLimit,Method,Signature"
		assert.EqualError(t, (&SignedMessage{}).Bind(req), expErr)
	})
}
