package v1_test

import (
	"context"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultHandler_ServeHTTP(t *testing.T) {
	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		port := test.RequireGetFreePort(t)
		s := server.NewHTTPAPI(context.Background(),
			&server.V1Callbacks{},
			port).
			Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		test.AssertResponseBody(t, port, "control/node", "/api/filecoin/v1/control/node is not implemented")
	})
}
