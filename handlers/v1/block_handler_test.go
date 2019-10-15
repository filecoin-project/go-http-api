package v1_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestBlockHeaderHandler_ServeHTTP(t *testing.T) {

	t.Run("all fields pass through", func(t *testing.T) {
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
		tb := types.Block{
			ID:     test.RequireTestCID(t, []byte("block")),
			Header: tbh,
		}

		bhh := v1.BlockHandler{Callback: func(blockId string) (*types.Block, error) {
			return &tb, nil
		}}

		cbs := &v1.Callbacks{GetBlockByID: bhh.Callback}
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
		assert.Equal(t, tb, actual)
	})

}
