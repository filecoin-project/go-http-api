package test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

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
	resp, err := http.Get(uri)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	return body
}

func AssertResponseBody(t *testing.T, port int, path string, exp string) {
	body := RequireGetResponseBody(t, port, path)
	assert.Equal(t, exp, string(body[:]))
}

func RequireTestCID(t *testing.T, data []byte) cid.Cid {
	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	require.NoError(t, err)
	return cid.NewCidV1(cid.DagCBOR, hash)
}
