package v1_test

import (
	"testing"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
)

func TestHelloHandler_ServeHTTP(t *testing.T) {
	t.Run("basic hello returns good response", func(t *testing.T) {
		cbs := &v1.Callbacks{}
		test.AssertServerResponse(t, cbs, false, "hello", "HELLO")
	})
}
