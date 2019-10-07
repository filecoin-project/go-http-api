package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/magiconair/properties/assert"
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

		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{GetActorByID: acb},
			port).
			Run()

		body := test.RequireGetResponseBody(t, port, "actors/1111")
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

		fmt.Println(string(errs[:]))

		acb := func(actorId string) (*types.Actor, error) {
			return nil, errors.New("this is an error")
		}

		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{GetActorByID: acb},
			port).
			Run()

		test.AssertResponseBody(t, port, "actors/1111", string(errs[:]))

	})

}
