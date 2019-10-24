package types

import (
	"encoding/json"
	"math/big"
	"net/http"
)

// Message is a struct representing a Filecoin message.
// Parameter are parameters for the Method from the HTTP request and are passed to the Callback unparsed.
type Message struct {
	Kind       string   `json:"kind"`
	ID         string   `json:"id,omitempty"`
	Nonce      uint64   `json:"nonce,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	Value      *big.Int `json:"value,omitempty"`    // in AttoFIL
	GasPrice   *big.Int `json:"gasPrice,omitempty"` // in AttoFIL
	GasLimit   uint64   `json:"gasLimit,omitempty"` // in GasUnits
	Method     string   `json:"method,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

// MarshalJSON marshals a Message struct
func (m Message) MarshalJSON() ([]byte, error) {
	type alias Message
	out := alias(m)
	out.Kind = "message"
	return json.Marshal(out)
}

// Bind does any additional setup needed for binding an http request body to
// a Message object
func (m *Message) Bind(r *http.Request) error {
	return RequireFields(m, "To", "Value", "GasPrice", "GasLimit", "Method")
}
