package server_test

import (
	"context"
	"testing"

	server "github.com/filecoin-project/go-http-api"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestNewHTTPServer(t *testing.T) {
	t.Run("if port is <=0 the default of :8080 will be used.", func(t *testing.T) {
		s := server.NewHTTPAPI(context.Background(), &v1.Callbacks{}, server.Config{}).Run()
		defer func() {
			s.Shutdown() // nolint: errcheck
		}()

		test.AssertGetResponseBodyEquals(t, 8080, false, "hello", "HELLO")
	})
}

func TestHTTPServer_Run(t *testing.T) {
	t.Run("basic hello returns good response", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, false, "hello", "HELLO")
	})

	t.Run("HTTPS requests", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, true, "hello", "HELLO")
	})

	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		nidcb := func() (*types.Node, error) {
			return &types.Node{}, nil
		}

		cbs := &v1.Callbacks{GetNode: nidcb}
		test.AssertServerResponse(t, cbs, false, "control/node", `{"kind":"node","protocol":{},"bitswapStats":{}}`)
	})

	t.Run("returns 404 when a path does not match", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, false, "foo", "404 page not found\n")
	})

	t.Run("if a handler was not provided a callback, returns 'is not implemented'", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, false, "control/node", "/api/filecoin/v1/control/node is not implemented")
	})
}
