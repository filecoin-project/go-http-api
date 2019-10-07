package server_test

import (
	"context"
	"testing"

	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"

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
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, server.Config{}).Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, 8080, "hello", "/api/filecoin/v1/hello, world!")
	})
}

func TestHTTPServer_Run(t *testing.T) {
	apiConfig := func() server.Config {
		return server.Config{
			Port: test.RequireGetFreePort(t),
		}
	}

	t.Run("basic hello returns good response", func(t *testing.T) {
		config := apiConfig()
		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			config).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, config.Port, "hello", "/api/filecoin/v1/hello, world!")
	})

	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		config := apiConfig()

		exp := `{"node":"node","protocol":{},"bitswapStats":{}}`
		nidcb := func() (*types.Node, error) {
			return &types.Node{}, nil
		}

		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{GetNode: nidcb},
			config).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, config.Port, "control/node", exp)
	})

	t.Run("returns 404 when a path does not match", func(t *testing.T) {
		config := apiConfig()

		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			config).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		path := "doesn't matter"

		test.AssertResponseBody(t, config.Port, path, "404 page not found\n")
	})
}
