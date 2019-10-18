package types

import (
	"math/big"
	"net/http"
)

// SignedMessage is a struct representing a signed message to be posted in the message pool.
// Unlike filecoin implementations, this does not embed Message and add Sigature,
// because of the need  to marshal as JSON.  Embedding Message causes json.Marshal to strip the Signature field.
// This is unlikely to be fixed: https://github.com/golang/go/issues/31167
type SignedMessage struct {
	Kind       string   `json:"kind,omitempty"`
	ID         string   `json:"id,omitempty"`
	Nonce      uint64   `json:"nonce,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	Value      *big.Int `json:"value,omitempty"`    // in AttoFIL
	GasPrice   *big.Int `json:"gasPrice,omitempty"` // in AttoFIL
	GasLimit   uint64   `json:"gasLimit,omitempty"` // in GasUnits
	Method     string   `json:"method,omitempty"`
	Parameters []string `json:"parameters"`
	Signature  string   `json:"signature,omitempty"`
}

func (sm *SignedMessage) Bind(r *http.Request) error {
	return RequireFields(sm, "To", "From", "Value", "GasPrice", "GasLimit", "Method", "Signature")
}
