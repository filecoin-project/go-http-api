package server_test

import (
	"context"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

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
		require.NoError(t, err)
		assert.Equal(t, "/hello, world!", string(body[:]))
	})
}

func TestHTTPServer_Run(t *testing.T) {
	t.Run("basic hello returns good response", func(t *testing.T) {
		port, err := test.GetFreePort()
		require.NoError(t, err)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()

		exp := []byte("/hello, world!")
		test.AssertResponseBody(t, port, "/hello", exp)
	})

	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()

		exp := []byte("/control/node is not implemented")
		test.AssertResponseBody(t, port, "/control/node", exp)
	})

	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		exp := []byte("abcd123")

		nidcb := func() ([]byte, error) {
			return exp, nil
		}

		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Node: nidcb},
			port).
			Run()

		test.AssertResponseBody(t, port, "/control/node", exp)
	})

	t.Run("returns 404 when a path does not match required param", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		acb := func(actorId string) ([]byte, error) {
			resp := []byte("doesn't matter")
			return resp, nil
		}

		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Actor: acb},
			port).
			Run()

		test.AssertResponseBody(t, port, "actors/", []byte("404 page not found\n"))
	})
}
