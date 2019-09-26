package server_test

import (
	"context"
	"fmt"
	server "github.com/carbonfive/go-filecoin-rest-api"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPServer_Hello(t *testing.T) {
	port, err := test.GetFreePort()
	require.NoError(t, err)
	s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, port)
	s.Run()

	uri := fmt.Sprintf("http://localhost:%d/hello", port)
	resp, err := http.Get(uri)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "/hello, world!", string(body[:]))
}

func TestNewHTTPServer(t *testing.T) {
	t.Run("if port is <=0 the default of :8080 will be used.", func(t *testing.T) {
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, 0)
		s.Run()

		resp, err := http.Get("http://localhost:8080/hello")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "/hello, world!", string(body[:]))
	})
}

func TestHTTPServer_Run(t *testing.T) {
	t.Run("calls default handler if no callback was provided", func(t *testing.T) {
		port, err := test.GetFreePort()
		require.NoError(t, err)
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{}, port)
		s.Run()

		uri := fmt.Sprintf("http://localhost:%d/node_id", port)
		resp, err := http.Get(uri)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "/node_id is not implemented", string(body[:]))
	})
	t.Run("calls correct handler if a callback for it was provided", func(t *testing.T) {
		nidcb := func() (string, error) {
			return "1234abcd", nil
		}

		port, err := test.GetFreePort()
		require.NoError(t, err)
		s := server.NewHTTPAPI(context.Background(), &server.V1Callbacks{NodeID: nidcb}, port)
		s.Run()

		uri := fmt.Sprintf("http://localhost:%d/node_id", port)
		resp, err := http.Get(uri)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "1234abcd", string(body[:]))
	})
}

func TestReflections(t *testing.T) {
	type foo struct {
		bar func(int)
		bazz func(string)
	}

}
