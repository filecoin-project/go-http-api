package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
)

type Message struct {
	Kind       string   `json:"kind,required,omitempty"`
	ID         string   `json:"id,omitempty"`
	Nonce      uint64   `json:"nonce,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	Value      *big.Int `json:"value,omitempty"`    // in AttoFIL
	GasPrice   *big.Int `json:"gasPrice,omitempty"` // in AttoFIL
	GasLimit   uint64   `json:"gasLimit,omitempty"` // in GasUnits
	Method     string   `json:"method,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
	Signature  string   `json:"signature,omitempty"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	type alias Message
	out := alias(m)
	out.Kind = "message"
	return json.Marshal(out)
}

type MessageRequest struct {
	*Message
}

func (mr *MessageRequest) Bind(r *http.Request) error {
	if mr.Message == nil {
		return errors.New("message fields missing")
	}
	return nil
}
