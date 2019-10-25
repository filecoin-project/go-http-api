package v1_test

import (
	"errors"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNodeHandler_ServeHTTP(t *testing.T) {
	testURL := "http://localhost:5000/control/node"

	t.Run("returns Node when callback succeeds", func(t *testing.T) {
		h := &v1.NodeHandler{Callback: happyNodeCb}
		rr := test.GetTestRequest(testURL, nil, h)
		assert.Equal(t, http.StatusOK, rr.Code)

		expBody := `{"kind":"node","id":"abcd123","protocol":{},"bitswapStats":{}}`
		assert.Equal(t, expBody, rr.Body.String())

		// lazytest: check that it works through server too
		cbs := &v1.Callbacks{GetNode: happyNodeCb}
		test.AssertServerResponse(t, cbs, false, "control/node", expBody)
	})

	t.Run("returns error when callback fails", func(t *testing.T) {
		h := &v1.NodeHandler{Callback: sadNodeCb}
		rr := test.GetTestRequest(testURL, nil, h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expBody := `{"errors":["boom"]}`
		assert.Equal(t, expBody, rr.Body.String())
	})
}

func happyNodeCb() (*types.Node, error) {
	return &types.Node{
		Kind:         "node",
		ID:           "abcd123",
		Protocol:     types.Protocol{},
		BitswapStats: types.BitswapStats{},
	}, nil
}

func sadNodeCb() (*types.Node, error) {
	return nil, errors.New("boom")
}
