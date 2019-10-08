package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/carbonfive/go-filecoin-rest-api/types"
)

func TestActor_ServeHTTP(t *testing.T) {
	t.Run("GetActorByID accepts and uses param actorId", func(t *testing.T) {
		fa := types.Actor{
			ActorType: "account",
			Address:   "abcd123",
			Balance:   *big.NewInt(600),
			Code:      test.RequireTestCID(t, []byte("anything")),
			Nonce:     123434,
		}
		acb := func(actorId string) (*types.Actor, error) {
			return &fa, nil
		}

		s := test.CreateTestServer(t, &server.V1Callbacks{GetActorByID: acb}, false)
		s.Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "actors/1111")
		var actual types.Actor
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, "actor", actual.Kind)
		assert.Equal(t, fa.ActorType, actual.ActorType)
		assert.Equal(t, fa.Address, actual.Address)
		assert.Equal(t, fa.Balance, actual.Balance)
		assert.Equal(t, fa.Code, actual.Code)
		assert.Equal(t, fa.Nonce, actual.Nonce)
	})

	t.Run("Errors are put into errors array", func(t *testing.T) {
		errs := types.MarshalErrors([]string{"this is an error"})

		acb := func(actorId string) (*types.Actor, error) {
			return nil, errors.New("this is an error")
		}

		cbs := &server.V1Callbacks{GetActorByID: acb}
		test.AssertServerResponse(t, cbs, false, "actors/1111", string(errs[:]))
	})
}
