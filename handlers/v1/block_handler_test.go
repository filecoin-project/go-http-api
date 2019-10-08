package v1_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	server "github.com/carbonfive/go-filecoin-rest-api"
	v1 "github.com/carbonfive/go-filecoin-rest-api/handlers/v1"
	"github.com/carbonfive/go-filecoin-rest-api/test"
	"github.com/carbonfive/go-filecoin-rest-api/types"
)

func TestBlockHeaderHandler_ServeHTTP(t *testing.T) {

	t.Run("all fields pass through", func(t *testing.T) {
		tbh := types.BlockHeader{
			Miner:                 "someaddress",
			Tickets:               [][]byte{[]byte("ticket1")},
			ElectionProof:         []byte("electionproof"),
			Parents:               []cid.Cid{test.RequireTestCID(t, []byte("parent1"))},
			ParentWeight:          *big.NewInt(1234),
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

		cbs := &server.V1Callbacks{GetBlockByID: bhh.Callback}
		s := test.CreateTestServer(t, cbs, false)
		s.Run()
		defer func() {
			assert.NoError(t, s.Shutdown())
		}()

		body := test.RequireGetResponseBody(t, s.Config().Port, "chain/blocks/1111")
		var actual types.Block
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.True(t, actual.ID.Equals(tb.ID))
		assert.Equal(t, "block", tb.Kind)
		assert.Equal(t, "blockHeader", tb.Header.Kind)
		assert.Equal(t, actual.Header.Miner, tb.Header.Miner)
		assert.Equal(t, actual.Header.Tickets[0], tb.Header.Tickets[0])
		assert.Equal(t, actual.Header.ElectionProof, tb.Header.ElectionProof)
		assert.True(t, tb.Header.Parents[0].Equals(actual.Header.Parents[0]))
		assert.Equal(t, actual.Header.ParentWeight, tb.Header.ParentWeight)
		assert.Equal(t, actual.Header.Height, tb.Header.Height)
		assert.True(t, tb.Header.ParentStateRoot.Equals(actual.Header.ParentStateRoot))
		assert.True(t, tb.Header.ParentMessageReceipts.Equals(actual.Header.ParentMessageReceipts))
		assert.True(t, tb.Header.Messages.Equals(actual.Header.Messages))
		assert.Equal(t, actual.Header.BLSAggregate, tb.Header.BLSAggregate)
		assert.Equal(t, actual.Header.Timestamp, tb.Header.Timestamp)
		assert.Equal(t, actual.Header.BlockSig, tb.Header.BlockSig)
	})

}
