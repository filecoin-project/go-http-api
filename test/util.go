package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"

	server "github.com/filecoin-project/go-http-api"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GetFreePort gets a free port from the kernel
// Credit: https://github.com/phayes/freeport
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() // nolint: errcheck
	return l.Addr().(*net.TCPAddr).Port, nil
}

func RequireGetFreePort(t *testing.T) int {
	port, err := GetFreePort()
	require.NoError(t, err)
	return port
}

func RequireGetResponseBody(t *testing.T, port int, path string) []byte {
	uri := fmt.Sprintf("http://localhost:%d/api/filecoin/v1/%s", port, path)
	return getResponseBody(t, uri)
}

func RequirePostFormResponseBody(t *testing.T, port int, path string, params url.Values) []byte {
	uri := fmt.Sprintf("http://localhost:%d/api/filecoin/v1/%s", port, path)
	return postFormResponseBody(t, uri, params)
}

func RequireGetResponseBodySSL(t *testing.T, port int, path string) []byte {
	uri := fmt.Sprintf("https://localhost:%d/api/filecoin/v1/%s", port, path)
	return getResponseBody(t, uri)
}

func AssertGetResponseBody(t *testing.T, port int, ssl bool, path string, exp string) {
	var body []byte

	if ssl {
		body = RequireGetResponseBodySSL(t, port, path)
	} else {
		body = RequireGetResponseBody(t, port, path)
	}
	assert.Equal(t, exp, string(body[:]))
}

func RequireTestCID(t *testing.T, data []byte) cid.Cid {
	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	require.NoError(t, err)
	return cid.NewCidV1(cid.DagCBOR, hash)
}

func getResponseBody(t *testing.T, uri string) []byte {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(uri)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	return body
}

func postFormResponseBody(t *testing.T, uri string, params url.Values) []byte {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	resp, err := client.PostForm(uri, params)
	require.NoError(t, err)
	require.Greater(t, 201, resp.StatusCode)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	return body
}

func CreateTestServer(t *testing.T, callbacks *server.V1Callbacks, ssl bool) *server.HTTPAPI {
	cfg := server.Config{Port: RequireGetFreePort(t)}
	if ssl {
		cfg.TLSCertPath = "test/fixtures/cert.pem"
		cfg.TLSKeyPath = "test/fixtures/key.pem"
	}

	return server.NewHTTPAPI(context.Background(), callbacks, cfg)
}

func AssertServerResponse(t *testing.T, callbacks *server.V1Callbacks, ssl bool, path string, expected string) {
	s := CreateTestServer(t, callbacks, ssl)

	s.Run()
	defer func() {
		assert.NoError(t, s.Shutdown())
	}()

	AssertGetResponseBody(t, s.Config().Port, ssl, path, expected)
}
