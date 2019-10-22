package types_test

import (
	"math/big"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/types"
)

func TestBlock_MarshalJSON(t *testing.T) {
	cid1 := test.RequireTestCID(t, []byte("cid1"))

	t.Run("struct is correctly serialized and includes Kind field", func(t *testing.T) {
		b := types.Block{
			ID: cid1,
			Header: types.BlockHeader{
				Miner:        "abcd",
				ParentWeight: big.NewInt(1),
				Height:       10,
				Timestamp:    38383838,
			},
		}
		expected := `{"kind":"block","ID":{"/":"bafyreib5znwh4i7pjxrtna4kzfuvhwnfklazpz6pe5ih4je2tyv7wmyesa"},"header":{"kind":"blockHeader","minerAddress":"abcd","parentWeight":1,"height":10,"parentStateRoot":null,"parentMessageReceipts":null,"messages":null,"timestamp":38383838}}`
		test.AssertMarshaledEquals(t, b, expected)
	})
	t.Run("empty struct is correctly serialized and includes Kind field", func(t *testing.T) {
		b := types.Block{}
		expected := `{"kind":"block","ID":null,"header":{"kind":"blockHeader","parentStateRoot":null,"parentMessageReceipts":null,"messages":null}}`
		test.AssertMarshaledEquals(t, b, expected)
	})
}

func TestBlockHeader_MarshalJSON(t *testing.T) {

	t.Run("full struct is correctly serialized and includes Kind field", func(t *testing.T) {
		cid1 := test.RequireTestCID(t, []byte("cid1"))
		cid2 := test.RequireTestCID(t, []byte("cid2"))
		cid3 := test.RequireTestCID(t, []byte("cid3"))
		cid4 := test.RequireTestCID(t, []byte("cid3"))

		bh := types.BlockHeader{
			Miner:                 "abcd",
			Tickets:               [][]byte{[]byte("abcd")},
			ElectionProof:         []byte("dcba"),
			Parents:               []cid.Cid{cid1},
			ParentWeight:          big.NewInt(1),
			Height:                10,
			ParentStateRoot:       cid2,
			ParentMessageReceipts: cid3,
			Messages:              cid4,
			BLSAggregate:          []byte("blsagg"),
			Timestamp:             38383838,
			BlockSig:              []byte("blocksig"),
		}

		expected := `{"kind":"blockHeader","minerAddress":"abcd","tickets":["YWJjZA=="],"electionProof":"ZGNiYQ==","parents":[{"/":"bafyreib5znwh4i7pjxrtna4kzfuvhwnfklazpz6pe5ih4je2tyv7wmyesa"}],"parentWeight":1,"height":10,"parentStateRoot":{"/":"bafyreihsffulhx7afspy7vmg3mo7nsau556h2kwlrtxjdvrpyg5iqgg33q"},"parentMessageReceipts":{"/":"bafyreiazoneogbvm4nsws53a2cbouhvb5ggaq7shpsoyjyyndkb4oft5wm"},"messages":{"/":"bafyreiazoneogbvm4nsws53a2cbouhvb5ggaq7shpsoyjyyndkb4oft5wm"},"blsAggregate":"YmxzYWdn","timestamp":38383838,"blockSig":"YmxvY2tzaWc="}`
		test.AssertMarshaledEquals(t, bh, expected)
	})

	t.Run("empty struct json omits all but special types and includes Kind field", func(t *testing.T) {
		emptyBh := types.BlockHeader{}

		expected := `{"kind":"blockHeader","parentStateRoot":null,"parentMessageReceipts":null,"messages":null}`
		test.AssertMarshaledEquals(t, emptyBh, expected)
	})
}
