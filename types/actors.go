package types

import (
	"math/big"

	"github.com/ipfs/go-cid"
)

type readableFunctionSignature struct {
	Params []string `json:"params"`
	Return []string `json:"return"`
}

// Actor is a struct for a Filecoin actor
type Actor struct {
	Kind      string                               `json:"kind,required,omitempty"`
	ActorType string                               `json:"actorType,omitempty"`
	Address   string                               `json:"address,omitempty"`
	Code      cid.Cid                              `json:"code,omitempty"`
	Nonce     uint64                               `json:"nonce,omitempty"`
	Balance   big.Int                              `json:"balance,omitempty"`
	Exports   map[string]readableFunctionSignature `json:"exports,omitempty"` // exports by function name
	Head      cid.Cid                              `json:"head,omitempty"`
}

type Actors struct {
	List []*Actor `json:"actors"`
}
