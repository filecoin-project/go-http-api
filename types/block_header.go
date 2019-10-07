package types

import (
	"math/big"

	"github.com/ipfs/go-cid"
)

// BlockHeader is the struct for a Filecoin blockheader

type Block struct {
	Kind   string `json:"kind,required"`
	ID     cid.Cid
	Header BlockHeader `json:"header"`
}

type BlockHeader struct {
	Kind                  string    `json:"kind,required,omitempty"`
	Miner                 string    `json:"minerAddress,omitempty"`
	Tickets               [][]byte  `json:"tickets,omitempty"`
	ElectionProof         []byte    `json:"electionProof,omitempty"`
	Parents               []cid.Cid `json:"parents,omitempty"`
	ParentWeight          big.Int   `json:"parentWeight,omitempty"`
	Height                uint64    `json:"height,omitempty"`
	ParentStateRoot       cid.Cid   `json:"parentStateRoot,omitempty"`
	ParentMessageReceipts cid.Cid   `json:"parentMessageReceipts,omitempty"`
	Messages              cid.Cid   `json:"messages,omitempty"`
	BLSAggregate          []byte    `json:"blsAggregate,omitempty"`
	Timestamp             uint64    `json:"timestamp,omitempty"`
	BlockSig              []byte    `json:"blockSig,omitempty"`
}
