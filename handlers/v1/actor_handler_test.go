package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
)

func TestActor_ServeHTTP(t *testing.T) {
	uri := "http://localhost/actors/1111"

	t.Run("GetActorByID accepts and uses param actorId", func(t *testing.T) {
		fa := requireCreateTestActor(t, fixtures.TestAddress0)
		h := &v1.ActorHandler{Callback: func(actorId string) (*types.Actor, error) {
			return fa, nil
		}}

		rr := test.GetTestRequest(uri, url.Values{}, h)
		assert.Equal(t, http.StatusOK, rr.Code)

		fa.Kind = "actor"
		var actual types.Actor
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual))
		assert.Equal(t, *fa, actual)
	})

	t.Run("Callback errors are returned w/ server error status", func(t *testing.T) {
		errs := types.MarshalErrorStrings("this is an error")

		acb := func(actorId string) (*types.Actor, error) {
			return nil, errors.New("this is an error")
		}
		h := &v1.ActorHandler{Callback: acb}

		rr := test.GetTestRequest(uri, url.Values{}, h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, errs, rr.Body.Bytes())
	})
}

func TestActorHandler_Integration(t *testing.T) {
	t.Run("GetActorByID accepts and uses param actorId", func(t *testing.T) {
		fa := requireCreateTestActor(t, fixtures.TestAddress0)
		acb := func(actorId string) (*types.Actor, error) {
			return fa, nil
		}

		s := test.CreateTestServer(t, &v1.Callbacks{GetActorByID: acb}, false).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "actors/1111")
		fa.Kind = "actor"
		var actual types.Actor
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, *fa, actual)
	})
}

func requireCreateTestActor(t *testing.T, addr string) *types.Actor {
	return &types.Actor{
		ActorType: "account",
		Address:   addr,
		Balance:   big.NewInt(600),
		Code:      test.RequireTestCID(t, []byte("anything")),
		Nonce:     big.NewInt(123434),
	}
}
