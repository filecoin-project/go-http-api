package types

import (
	"encoding/json"
	"math/big"
	"net/http"
)

// Message is a struct representing a Filecoin message.
// Parameter are parameters for the Method from the HTTP request and are passed to the Callback unparsed.
type Message struct {
	Kind       string   `json:"kind,required,omitempty"`
	ID         string   `json:"id,omitempty"`
	Nonce      uint64   `json:"nonce,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,required,omitempty"`
	Value      *big.Int `json:"value,required,omitempty"`    // in AttoFIL
	GasPrice   *big.Int `json:"gasPrice,required,omitempty"` // in AttoFIL
	GasLimit   uint64   `json:"gasLimit,required,omitempty"` // in GasUnits
	Method     string   `json:"method,required,omitempty"`
	Parameters []string `json:"parameters,required"`
	Signature  string   `json:"signature,omitempty"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	type alias Message
	out := alias(m)
	out.Kind = "message"
	return json.Marshal(out)
}

func (m *Message) Bind(r *http.Request) error {
	return RequireFields(m, "To", "Value", "GasPrice", "GasLimit", "Method")
}
