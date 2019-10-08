package v1_test

import (
	"testing"

	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
)

func TestHelloHandler_ServeHTTP(t *testing.T) {
	t.Run("basic hello returns good response", func(t *testing.T) {
		cbs := &server.V1Callbacks{}
		test.AssertServerResponse(t, cbs, false, "hello", "/api/filecoin/v1/hello, world!")
	})
}
