package v1_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
)

func TestBlockHandler_ServeHTTP(t *testing.T) {
	uri := "http://chain/blocks/1111"
	t.Run("returns a block and status ok", func(t *testing.T) {
		tb := requireCreateTestBlock(t, fixtures.TestAddress1)
		h := &v1.BlockHandler{Callback: func(blockId string) (*types.Block, error) {
			return tb, nil
		}}
		rr := test.GetTestRequest(uri, url.Values{}, h)

		assert.Equal(t, http.StatusOK, rr.Code)

		tb.Kind = "block"
		expected, _ := json.Marshal(*tb)
		assert.Equal(t, expected, rr.Body.Bytes())
	})

	t.Run("if callback errors, returns server error status + error msg", func(t *testing.T) {
		err := errors.New("boom")
		h := &v1.BlockHandler{Callback: func(blockId string) (*types.Block, error) {
			return &types.Block{}, err
		}}
		rr := test.GetTestRequest(uri, url.Values{}, h)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		expected := types.MarshalError(err)
		assert.Equal(t, expected, rr.Body.Bytes())
	})
}

func TestBlockHeaderHandler_Integration(t *testing.T) {
	t.Run("all fields pass through", func(t *testing.T) {
		tb := requireCreateTestBlock(t, fixtures.TestAddress0)

		cbs := &v1.Callbacks{
			GetBlockByID: func(blockId string) (*types.Block, error) {
				return tb, nil
			}}
		s := test.CreateTestServer(t, cbs, false)
		s.Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "chain/blocks/1111")
		var actual types.Block
		require.NoError(t, json.Unmarshal(body, &actual))
		tb.Header.Kind = "blockHeader"
		tb.Kind = "block"
		assert.Equal(t, *tb, actual)
	})
}

func requireCreateTestBlock(t *testing.T, addr string) *types.Block {
	tbh := types.BlockHeader{
		Miner:                 "someaddress",
		Tickets:               [][]byte{[]byte("ticket1")},
		ElectionProof:         []byte("electionproof"),
		Parents:               []cid.Cid{test.RequireTestCID(t, []byte("parent1"))},
		ParentWeight:          big.NewInt(1234),
		Height:                34343,
		ParentStateRoot:       test.RequireTestCID(t, []byte("stateroot")),
		ParentMessageReceipts: test.RequireTestCID(t, []byte("receipts")),
		Messages:              test.RequireTestCID(t, []byte("messages")),
		BLSAggregate:          []byte("blsa"),
		Timestamp:             939393,
		BlockSig:              []byte("blocksig"),
	}
	return &types.Block{
		ID:     test.RequireTestCID(t, []byte("block")),
		Header: tbh,
	}
}
