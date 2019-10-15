package v1_test

import (
	"encoding/json"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

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

		req := httptest.NewRequest("GET", path, nil)
		rr := httptest.NewRecorder()
		handler := v1.ActorsHandler{Callback: acb}
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var actual []types.Actor
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))

		require.Len(t, actual, 2)
		assert.Equal(t, a1, actual[0])
		assert.Equal(t, a2, actual[1])
	})
}
