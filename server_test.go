package server_test

import (
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPServer_Hello(t *testing.T) {
	hello := &handlers.HelloHandler{}
	s := server.NewHTTPServer(hello)
	go s.Run()
	resp, err := http.Get("http://localhost:8080/hello")
	assert.NoError(t, err)
	if err != nil {
		// handle error
	}
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "/hello, world!", string(body[:]))
}