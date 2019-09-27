package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/handlers/api_errors"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"math/big"
	"testing"
)

func TestActor_ServeHTTP(t *testing.T) {
	t.Run("Actor accepts and uses param actorId", func(t *testing.T) {
		type FakeActor struct {
			ActorType string
			Address   string
			Nonce     uint64
			Balance   big.Int
		}

		fa := FakeActor{
			ActorType: "account",
			Address:   "abcd123",
			Nonce:     123434,
			Balance:   *big.NewInt(600),
		}
		fajson, _ := json.Marshal(fa)
		acb := func(actorId string) ([]byte, error) {
			resp, _ := json.Marshal(fa)
			return resp, nil
		}

		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Actor: acb},
			port).
			Run()

		test.AssertResponseBody(t, port, "actors/1111", fajson)
	})

	t.Run("Errors are put into errors array", func(t *testing.T) {
		errs := api_errors.MarshalErrors([]string{"this is an error"})

		fmt.Println(string(errs[:]))

		acb := func(actorId string) ([]byte, error) {
			return []byte{}, errors.New("this is an error")
		}

		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Actor: acb},
			port).
			Run()

		test.AssertResponseBody(t, port, "actors/1111", errs)

	})

}
