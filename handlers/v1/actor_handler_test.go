package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestActor_ServeHTTP(t *testing.T) {
	t.Run("GetActorByID accepts and uses param actorId", func(t *testing.T) {
		fa := types.Actor{
			ActorType: "account",
			Address:   "abcd123",
			Balance:   big.NewInt(600),
			Code:      test.RequireTestCID(t, []byte("anything")),
			Nonce:     123434,
		}
		acb := func(actorId string) (*types.Actor, error) {
			return &fa, nil
		}

		s := test.CreateTestServer(t, &v1.Callbacks{GetActorByID: acb}, false).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "actors/1111")
		fa.Kind = "actor"
		var actual types.Actor
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, fa, actual)
	})

	t.Run("Errors are put into errors array", func(t *testing.T) {
		errs := types.MarshalErrors([]string{"this is an error"})

		acb := func(actorId string) (*types.Actor, error) {
			return nil, errors.New("this is an error")
		}

		cbs := &v1.Callbacks{GetActorByID: acb}
		test.AssertServerResponse(t, cbs, false, "actors/1111", string(errs[:]))
	})
}
