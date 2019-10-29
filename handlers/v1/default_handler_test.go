package v1_test

import (
	"testing"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
)

func TestDefaultHandler_ServeHTTP(t *testing.T) {
	t.Run("calls default handler if no callback is provided", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, false, "control/node", "/api/filecoin/v1/control/node is not implemented")
	})
	t.Run("does not call default handler when callback is provided", func(t *testing.T) {
		cbs := &v1.Callbacks{GetActorNonce: happyPathANCallback}
		test.AssertServerResponse(t, cbs, false, "actors/1234/nonce", "54321")
	})
}
