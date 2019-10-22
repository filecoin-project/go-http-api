package types_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/filecoin-project/go-http-api/types"
)

func TestRequireFields(t *testing.T) {
	type Foo struct {
		AString    string
		AUint64    uint64
		ABigIntPtr *big.Int
		ASlice     []byte
	}

	t.Run("returns nil if all required fields are present", func(t *testing.T) {
		f := Foo{
			AString:    "foo",
			AUint64:    1234,
			ABigIntPtr: big.NewInt(1234),
			ASlice:     []byte("byte"),
		}
		assert.NoError(t, RequireFields(f, "AString", "AUint64", "ABigIntPtr", "ASlice"))
	})
	t.Run("Returns error if field doesn't exist in struct", func(t *testing.T) {
		f := Foo{
			AString:    "foo",
			AUint64:    1234,
			ABigIntPtr: big.NewInt(1234),
			ASlice:     []byte("byte"),
		}
		assert.EqualError(t, RequireFields(f, "AFieldThatDoesNotExist"), "AFieldThatDoesNotExist is not part of struct Foo")
	})
	t.Run("returns list of missing fields in error if fields are missing", func(t *testing.T) {
		f := Foo{}
		assert.EqualError(t, RequireFields(f, "AString", "AUint64", "ABigIntPtr", "ASlice"), "missing parameters: AString,AUint64,ABigIntPtr,ASlice")

	})
}
