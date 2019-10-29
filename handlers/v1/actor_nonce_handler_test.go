package v1_test

import (
	"errors"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestActorNonceHandler_ServeHTTP(t *testing.T) {
	uri := "http://localhost:3000/actors/1234/nonce"

	t.Run("returns actor nonce when requested", func(t *testing.T) {
		h := &v1.ActorNonceHandler{Callback: happyPathANCallback}
		rr := test.GetTestRequest(uri, nil, h)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "54321", rr.Body.String())
	})

	t.Run("returns error if callback fails", func(t *testing.T) {
		h := &v1.ActorNonceHandler{Callback: sadPathANCallback}
		rr := test.GetTestRequest(uri, nil, h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		exp := types.MarshalErrorStrings("boom")
		assert.Equal(t, exp, rr.Body.Bytes())
	})

	t.Run("returns actor nonce via server", func(t *testing.T) {
		cbs := &v1.Callbacks{GetActorNonce: happyPathANCallback}
		test.AssertServerResponse(t, cbs, false, "actors/1234/nonce", "54321")
	})
}

func happyPathANCallback(_ string) (*big.Int, error) {
	return big.NewInt(54321), nil
}

func sadPathANCallback(_ string) (*big.Int, error) {
	return nil, errors.New("boom")
}
