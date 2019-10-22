package types_test

import (
	"math/big"
	"testing"

	"github.com/filecoin-project/go-http-api/test"
	. "github.com/filecoin-project/go-http-api/types"
)

func TestActor_MarshalJSON(t *testing.T) {
	t.Run("full struct returns correct values plus Kind", func(t *testing.T) {
		rfs1 := ReadableFunctionSignature{
			Params: []string{},
			Return: []string{"string"},
		}

		exports := map[string]ReadableFunctionSignature{}
		exports["getNetwork"] = rfs1
		cid1 := test.RequireTestCID(t, []byte("cid1"))
		cid2 := test.RequireTestCID(t, []byte("cid2"))

		a := Actor{
			ActorType: "account",
			Address:   "abcd345",
			Code:      cid1,
			Nonce:     big.NewInt(1),
			Balance:   big.NewInt(2),
			Exports:   exports,
			Head:      cid2,
		}
		expected := `{"kind":"actor","role":"account","address":"abcd345","code":{"/":"bafyreib5znwh4i7pjxrtna4kzfuvhwnfklazpz6pe5ih4je2tyv7wmyesa"},"nonce":1,"balance":2,"exports":{"getNetwork":{"params":[],"return":["string"]}},"head":{"/":"bafyreihsffulhx7afspy7vmg3mo7nsau556h2kwlrtxjdvrpyg5iqgg33q"}}`
		test.AssertMarshaledEquals(t, a, expected)
	})
	t.Run("empty struct returns correct valudes plus Kind", func(t *testing.T) {
		a := Actor{}
		expected := `{"kind":"actor","code":null,"exports":null,"head":null}`
		test.AssertMarshaledEquals(t, a, expected)
	})
}
