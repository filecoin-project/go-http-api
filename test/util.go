package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	server "github.com/filecoin-project/go-http-api"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
)

type Param struct {
	Key   string
	Value string
}

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

func CreateTestServer(t *testing.T, callbacks *v1.Callbacks, ssl bool) *server.HTTPAPI {
	cfg := server.Config{Port: RequireGetFreePort(t)}
	if ssl {
		cfg.TLSCertPath = "test/fixtures/cert.pem"
		cfg.TLSKeyPath = "test/fixtures/key.pem"
	}

	return server.NewHTTPAPI(context.Background(), callbacks, cfg)
}

func AssertServerResponse(t *testing.T, callbacks *v1.Callbacks, ssl bool, path string, expected string) {
	s := CreateTestServer(t, callbacks, ssl)

	s.Run()
	defer func() {
		assert.NoError(t, s.Shutdown())
	}()

	AssertGetResponseBody(t, s.Config().Port, ssl, path, expected)
}

func GetTestRequest(uri string, params *[]Param, h http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", uri, nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	if params != nil {
		for _, el := range *params {
			rctx.URLParams.Add(el.Key, el.Value)
		}
	}

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.ServeHTTP(rr, req)
	return rr
}

func PostTestRequest(uri string, body io.Reader, h http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/json")
	req.PostForm = url.Values{}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
