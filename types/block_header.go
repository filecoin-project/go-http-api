package types

import (
	"math/big"

	"github.com/ipfs/go-cid"
)

// BlockHeader is the struct for a Filecoin blockheader
type BlockHeader struct {
	Miner                 string    `json:"minerAddress"`
	Tickets               [][]byte  `json:"tickets"`
	ElectionProof         []byte    `json:"electionProof"`
	Parents               []cid.Cid `json:"parents"`
	ParentWeight          big.Int   `json:"parentWeight"`
	Height                uint64    `json:"height"`
	ParentStateRoot       cid.Cid   `json:"parentStateRoot"`
	ParentMessageReceipts cid.Cid   `json:"parentMessageReceipts"`
	Messages              cid.Cid   `json:"messages"`
	BLSAggregate          []byte    `json:"blsAggregate"`
	Timestamp             uint64    `json:"timestamp"`
	BlockSig              []byte    `json:"blockSig"`
}
