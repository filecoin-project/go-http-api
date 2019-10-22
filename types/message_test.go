package types_test

import (
	"github.com/filecoin-project/go-http-api/test"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http/httptest"
	"testing"

	"github.com/filecoin-project/go-http-api/test/fixtures"
	. "github.com/filecoin-project/go-http-api/types"
)

func TestMessage_MarshalJSON(t *testing.T) {

	t.Run("full struct is correctly serialized and includes Kind field", func(t *testing.T) {
		tm := Message{
			ID:         "someid",
			Nonce:      100,
			From:       fixtures.TestAddress0,
			To:         fixtures.TestAddress1,
			Value:      big.NewInt(1),
			GasPrice:   big.NewInt(2),
			GasLimit:   3,
			Method:     "cancel",
			Parameters: []string{"abcd123"},
		}

		expected := `{"kind":"message","id":"someid","nonce":100,"from":"t2gmpzificaunkf47tzkt377a6yllmcfj3g3qbyti","to":"t12cvsox5neub6y4vupgsogbrfaljiot4eaenkkyy","value":1,"gasPrice":2,"gasLimit":3,"method":"cancel","parameters":["abcd123"]}`
		test.AssertMarshaledEquals(t, tm, expected)
	})

	t.Run("empty struct is correctly serialized and includes Kind field", func(t *testing.T) {
		tm := Message{}
		expected := `{"kind":"message"}`
		test.AssertMarshaledEquals(t, tm, expected)
	})
}

func TestMessage_Bind(t *testing.T) {
	t.Run("returns nil if required fields are present", func(t *testing.T) {
		tm := Message{
			ID:         "someid",
			Nonce:      100,
			From:       fixtures.TestAddress0,
			To:         fixtures.TestAddress1,
			Value:      big.NewInt(1),
			GasPrice:   big.NewInt(2),
			GasLimit:   3,
			Method:     "cancel",
			Parameters: []string{"abcd123"},
		}
		req := httptest.NewRequest("GET", "/some/path", nil)
		assert.NoError(t, tm.Bind(req))
	})
	t.Run("returns error if fields are missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some/path", nil)
		expErr := "missing parameters: To,Value,GasPrice,GasLimit,Method"
		assert.EqualError(t, (&Message{}).Bind(req), expErr)
	})
}
