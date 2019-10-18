package types

import (
	"net/http"
)

// SignedMessage is a struct representing a signed message to be posted in the message pool.
type SignedMessage struct {
	Message
	Signature string `json:"signature,omitempty"`
}

func (sm *SignedMessage) Bind(r *http.Request) error {
	return RequireFields(sm, "To", "Value", "GasPrice", "GasLimit", "Method", "Parameters", "Signature")
}
