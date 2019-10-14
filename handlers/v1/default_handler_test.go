package v1_test

import (
	"testing"

	server "github.com/filecoin-project/go-http-api"
	"github.com/filecoin-project/go-http-api/test"
)

func TestDefaultHandler_ServeHTTP(t *testing.T) {
	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		cbs := &server.V1Callbacks{}
		test.AssertServerResponse(t, cbs, false, "control/node", "/api/filecoin/v1/control/node is not implemented")
	})
}
