package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	server "github.com/filecoin-project/go-http-api"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
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

// RequireGetFreePort fails the test if GetFreePort fails
func RequireGetFreePort(t *testing.T) int {
	port, err := GetFreePort()
	require.NoError(t, err)
	return port
}

// RequireGetResponseBody fails the test if getResponseBody fails, when posted
// to a test server
func RequireGetResponseBody(t *testing.T, port int, path string) []byte {
	uri := fmt.Sprintf("http://localhost:%d/api/filecoin/v1/%s", port, path)
	return getResponseBody(t, uri)
}

// RequireGetResponseBodySSL fails the test if getResponseBody fails, when posted
// to a test server
func RequireGetResponseBodySSL(t *testing.T, port int, path string) []byte {
	uri := fmt.Sprintf("https://localhost:%d/api/filecoin/v1/%s", port, path)
	return getResponseBody(t, uri)
}

// AssertGetResponseBodyEquals asserts that response body for a GET call using the provided
// arguments equals `exp`, when posted to a test server
func AssertGetResponseBodyEquals(t *testing.T, port int, ssl bool, path string, exp string) {
	var body []byte

	if ssl {
		body = RequireGetResponseBodySSL(t, port, path)
	} else {
		body = RequireGetResponseBody(t, port, path)
	}
	assert.Equal(t, exp, string(body[:]))
}

// RequireTestCID generates a new random cid.Cid
func RequireTestCID(t *testing.T, data []byte) cid.Cid {
	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	require.NoError(t, err)
	return cid.NewCidV1(cid.DagCBOR, hash)
}

// getResponseBody gets a response from a running http server and returns it as
// []byte
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

// CreateTestServer creates a real http server for testing
func CreateTestServer(t *testing.T, callbacks *v1.Callbacks, ssl bool) *server.HTTPAPI {
	cfg := server.Config{Port: RequireGetFreePort(t)}
	if ssl {
		cfg.TLSCertPath = "test/fixtures/cert.pem"
		cfg.TLSKeyPath = "test/fixtures/key.pem"
	}

	return server.NewHTTPAPI(context.Background(), callbacks, cfg)
}

// AssertServerResponse creates an http test server, sends a request and asserts that
// the response body is equal to `exp`
func AssertServerResponse(t *testing.T, callbacks *v1.Callbacks, ssl bool, path string, exp string) {
	s := CreateTestServer(t, callbacks, ssl)

	s.Run()
	defer s.Shutdown() // nolint: errcheck

	AssertGetResponseBodyEquals(t, s.Config().Port, ssl, path, exp)
}

// GetTestRequest sets up a request to uri with url params via httptest, calls the
// provided handler, and returns the new recorder with the response stored.
func GetTestRequest(getURL string, params url.Values, h http.Handler) *httptest.ResponseRecorder {
	rctx := chi.NewRouteContext()
	req := httptest.NewRequest("GET", getURL, nil)
	req.Form = params

	// have to add the params to the chi context; otherwise chi doesn't know about them.
	if params != nil {
		for k, v := range params {
			rctx.URLParams.Add(k, strings.Join(v, ","))
		}
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

// PostTestRequest sets up a post request to `uri` with a JSON `body`, calls the
// provided handler `h` and returns the new recorder with the response stored.
func PostTestRequest(uri string, body io.Reader, h http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/json")
	req.PostForm = url.Values{}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

// Marshalable is an interface that has a `MarshalJSON` func, for use by
// AssertMarshaledEquals
type Marshalable interface {
	MarshalJSON() ([]byte, error)
}

// AssertMarshaledEquals asserts that the marshaled version of the Marshalable `m` equals
// `exp`
func AssertMarshaledEquals(t *testing.T, m Marshalable, exp string) {
	marshaled, err := m.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, exp, string(marshaled[:]))
}

// AssertEqualFuncs compares two funcs for equality
func AssertEqualFuncs(t *testing.T, fn1, fn2 interface{}) {
	assert.Equal(t, funcPtrAsString(fn1), funcPtrAsString(fn2))
}

// funcPtrAsString gets the pointer value of the func
func funcPtrAsString(fn interface{}) string {
	res := fmt.Sprintf("%v", fn)
	return res
}
