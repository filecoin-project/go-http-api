package types

import (
	"github.com/ipfs/go-cid"
	"math/big"
)

type readableFunctionSignature struct {
	Params []string `json:"params"`
	Return []string `json:"return"`
}

// V1 Actor
type Actor struct {
	ActorType string                               `json:"actorType"`
	Address   string                               `json:"address"`
	Code      cid.Cid                              `json:"code,omitempty"`
	Nonce     uint64                               `json:"nonce"`
	Balance   big.Int                              `json:"balance"`
	Exports   map[string]readableFunctionSignature `json:"exports"` // exports by function name
	Head      cid.Cid                              `json:"head,omitempty"`
}

type Actors struct {
	List []*Actor `json:"actors"`
}
