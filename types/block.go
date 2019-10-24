package types

import (
	"encoding/json"
	"math/big"

	"github.com/ipfs/go-cid"
)

// Block is the internal struct for a Filecoin block
type Block struct {
	Kind   string `json:"kind"`
	ID     cid.Cid
	Header BlockHeader `json:"header"`
}

// MarshalJSON marshals a Block struct
func (o Block) MarshalJSON() ([]byte, error) {
	type alias Block
	out := alias(o)
	out.Kind = "block"
	return json.Marshal(out)
}

// BlockHeader is the internal struct for a Filecoin blockheader
type BlockHeader struct {
	Kind                  string    `json:"kind"`
	Miner                 string    `json:"minerAddress,omitempty"`
	Tickets               [][]byte  `json:"tickets,omitempty"`
	ElectionProof         []byte    `json:"electionProof,omitempty"`
	Parents               []cid.Cid `json:"parents,omitempty"`
	ParentWeight          *big.Int  `json:"parentWeight,omitempty"`
	Height                uint64    `json:"height,omitempty"`
	ParentStateRoot       cid.Cid   `json:"parentStateRoot,omitempty"`
	ParentMessageReceipts cid.Cid   `json:"parentMessageReceipts,omitempty"`
	Messages              cid.Cid   `json:"messages,omitempty"`
	BLSAggregate          []byte    `json:"blsAggregate,omitempty"`
	Timestamp             uint64    `json:"timestamp,omitempty"`
	BlockSig              []byte    `json:"blockSig,omitempty"`
}

// MarshalJSON marshals a BlockHeader struct
func (o BlockHeader) MarshalJSON() ([]byte, error) {
	type alias BlockHeader
	out := alias(o)
	out.Kind = "blockHeader"
	return json.Marshal(out)
}
