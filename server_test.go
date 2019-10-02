package server_test

import (
	"context"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/require"
	"testing"
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
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, 0)
		s.Run()

		test.AssertResponseBody(t, 8080, "hello", "/api/filecoin/v1/hello, world!")
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

		test.AssertResponseBody(t, port, "hello", "/api/filecoin/v1/hello, world!")
	})

	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)
		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()

		test.AssertResponseBody(t, port, "control/node", "/api/filecoin/v1/control/node is not implemented")
	})

	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		exp := "abcd123"
		nidcb := func() ([]byte, error) {
			return []byte(exp), nil
		}

		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Node: nidcb},
			port).
			Run()

		test.AssertResponseBody(t, port, "control/node", exp)
	})

	t.Run("returns 404 when a path does not match", func(t *testing.T) {
		port := test.RequireGetFreePort(t)

		acb := func(actorId string) ([]byte, error) {
			resp := []byte("doesn't matter")
			return resp, nil
		}

		server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{Actor: acb},
			port).
			Run()

		test.AssertResponseBody(t, port, "foo", "404 page not found\n")
	})
}
