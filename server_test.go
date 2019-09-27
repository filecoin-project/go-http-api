package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
)

func TestHTTPServer_Hello(t *testing.T) {
	port, err := test.GetFreePort()
	require.NoError(t, err)
	s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, port)
	s.Run()

	uri := fmt.Sprintf("http://localhost:%d/hello", port)
	resp, err := http.Get(uri)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "/hello, world!", string(body[:]))
}

func TestNewHTTPServer(t *testing.T) {
	t.Run("if port is <=0 the default of :8080 will be used.", func(t *testing.T) {
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, 0)
		s.Run()

		resp, err := http.Get("http://localhost:8080/hello")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "/hello, world!", string(body[:]))
	})
}

func TestHTTPServer_Run(t *testing.T) {
	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		port := requireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()

		uri := fmt.Sprintf("http://localhost:%d/control/node", port)
		resp, err := http.Get(uri)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "/control/node is not implemented", string(body[:]))
	})
	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		port := requireGetFreePort(t)

		exp := []byte("abcd123")

		nidcb := func() ([]byte, error) {
			return exp, nil
		}

		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Node: nidcb},
			port).
			Run()

		assertResponseBody(t, port, "/control/node", exp)
	})

	t.Run("Actor route accepts an id", func(t *testing.T) {
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

		port := requireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Actor: acb},
			port).
			Run()

		assertResponseBody(t, port, "actors/1111", fajson)
	})
}

func requireGetFreePort(t *testing.T) int {
	port, err := test.GetFreePort()
	require.NoError(t, err)
	return port
}

func assertResponseBody(t *testing.T, port int, path string, exp []byte) {
	uri := fmt.Sprintf("http://localhost:%d/%s", port, path)
	resp, err := http.Get(uri)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, exp, body)

}
