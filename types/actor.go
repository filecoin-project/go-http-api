package types

import (
	"encoding/json"
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
	ActorType string                               `json:"role,omitempty"`
	Address   string                               `json:"address,omitempty"`
	Code      cid.Cid                              `json:"code,omitempty"`
	Nonce     *big.Int                             `json:"nonce,omitempty"`
	Balance   *big.Int                             `json:"balance,omitempty"`
	Exports   map[string]readableFunctionSignature `json:"exports,omitempty"` // exports by function name
	Head      cid.Cid                              `json:"head,omitempty"`
}

func (o Actor) MarshalJSON() ([]byte, error) {
	type alias Actor
	out := alias(o)
	out.Kind = "actor"
	return json.Marshal(out)
}
