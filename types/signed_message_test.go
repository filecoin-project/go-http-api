package types_test

import (
	"github.com/filecoin-project/go-http-api/test"
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

func TestSignedMessage_MarshalJSON(t *testing.T) {
	t.Run("full struct is correctly serialized and includes Kind field", func(t *testing.T) {
		tm := &SignedMessage{
			ID:         "someid",
			Nonce:      100,
			From:       fixtures.TestAddress0,
			To:         fixtures.TestAddress1,
			Value:      big.NewInt(1),
			GasPrice:   big.NewInt(2),
			GasLimit:   3,
			Method:     "cancel",
			Parameters: []string{"abcd123"},
			Signature:  "somesig",
		}

		expected := `{"kind":"signedMessage","id":"someid","nonce":100,"from":"t2gmpzificaunkf47tzkt377a6yllmcfj3g3qbyti","to":"t12cvsox5neub6y4vupgsogbrfaljiot4eaenkkyy","value":1,"gasPrice":2,"gasLimit":3,"method":"cancel","parameters":["abcd123"],"signature":"somesig"}`
		test.AssertMarshaledEquals(t, tm, expected)
	})

	t.Run("empty struct is correctly serialized and includes Kind field", func(t *testing.T) {
		tm := &SignedMessage{}
		expected := `{"kind":"signedMessage","parameters":null}`
		test.AssertMarshaledEquals(t, tm, expected)

	})
}
