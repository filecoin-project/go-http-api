package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
)

func TestActorsHandler_ServeHTTP(t *testing.T) {
	path := "http://localhost/actors"
	t.Run("Returns list of actors associated with node", func(t *testing.T) {
		a1 := types.Actor{
			ActorType: "account",
			Address:   fixtures.TestAddress0,
			Balance:   big.NewInt(600),
			Code:      test.RequireTestCID(t, []byte("actor1")),
			Nonce:     21,
			Head:      test.RequireTestCID(t, []byte("head")),
		}
		a2 := types.Actor{
			ActorType: "miner",
			Address:   fixtures.TestAddress1,
			Code:      test.RequireTestCID(t, []byte("actor2")),
			Nonce:     10,
			Balance:   big.NewInt(2100),
			Head:      test.RequireTestCID(t, []byte("head1")),
		}

		acb := func() ([]*types.Actor, error) {
			a1.Kind = "actor"
			a2.Kind = "actor"
			return []*types.Actor{&a1, &a2}, nil
		}

		h := v1.ActorsHandler{Callback: acb}
		resp, body := test.TestGetHandler(&h, path, nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var actual []types.Actor
		require.NoError(t, json.Unmarshal(body, &actual))

		require.Len(t, actual, 2)
		assert.Equal(t, a1, actual[0])
		assert.Equal(t, a2, actual[1])
	})

	t.Run("Returns error message with server error if callback fails", func(t *testing.T) {
		expectedErrs := types.MarshalErrors([]string{"boom"})

		acb := func() ([]*types.Actor, error) {
			return []*types.Actor{}, errors.New("boom")
		}

		h := v1.ActorsHandler{Callback: acb}
		resp, body := test.TestGetHandler(&h, path, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, expectedErrs, body)
	})
}
