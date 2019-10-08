package v1_test

import (
	"context"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHelloHandler_ServeHTTP(t *testing.T) {
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
}
