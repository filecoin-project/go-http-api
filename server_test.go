package server_test

import (
	"context"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
)

//func TestRoutes(t *testing.T) {
//	port, err := test.GetFreePort()
//	require.NoError(t, err)
//	s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{
//
//	}, port)
//}

func TestNewHTTPServer(t *testing.T) {
	t.Run("if port is <=0 the default of :8080 will be used.", func(t *testing.T) {
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, 0).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, 8080, "hello", "/api/filecoin/v1/hello, world!")
	})
}

func TestHTTPServer_Run(t *testing.T) {
	t.Run("basic hello returns good response", func(t *testing.T) {
		port, err := test.GetFreePort()
		require.NoError(t, err)
		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, port, "hello", "/api/filecoin/v1/hello, world!")
	})

	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		exp := "{\"node\":\"node\",\"protocol\":{},\"bitswapStats\":{}}"

		nidcb := func() (*types.Node, error) {
			return &types.Node{}, nil
		}

		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{GetNode: nidcb},
			port).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, port, "control/node", exp)
	})

	t.Run("returns 404 when a path does not match", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		path := "doesn't matter"

		test.AssertResponseBody(t, port, path, "404 page not found\n")
	})
}
