package types

import (
	"encoding/json"
	"math/big"

	"github.com/ipfs/go-cid"
)

// ReadableFunctionSignature represents a function signature definition for exported
// Filecoin functions, i.e. those that can be invoked by a message send
type ReadableFunctionSignature struct {
	Params []string `json:"params"`
	Return []string `json:"return"`
}

// Actor is a struct for a Filecoin actor
type Actor struct {
	Kind      string                               `json:"kind"`
	ActorType string                               `json:"role,omitempty"`
	Address   string                               `json:"address,omitempty"`
	Code      cid.Cid                              `json:"code,omitempty"`
	Nonce     *big.Int                             `json:"nonce,omitempty"`
	Balance   *big.Int                             `json:"balance,omitempty"`
	Exports   map[string]ReadableFunctionSignature `json:"exports"` // exports by function name
	Head      cid.Cid                              `json:"head,omitempty"`
}

// MarshalJSON marshals Actor into JSON
func (o Actor) MarshalJSON() ([]byte, error) {
	type alias Actor
	out := alias(o)
	out.Kind = "actor"
	return json.Marshal(out)
}
