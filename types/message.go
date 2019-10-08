package types

import "math/big"

type Message struct {
	Kind       string        `json:"kind,required,omitempty"`
	ID         string        `json:"id,omitempty"`
	Nonce      uint64        `json:"nonce,omitempty"`
	From       string        `json:"from,omitempty"`
	To         string        `json:"to,omitempty"`
	Value      *big.Int      `json:"value,omitempty"`    // in AttoFIL
	GasPrice   *big.Int      `json:"gasPrice,omitempty"` // in AttoFIL
	GasLimit   uint64        `json:"gasLimit,omitempty"` // in GasUnits
	Method     string        `json:"method,omitempty"`
	Parameters []interface{} `json:"parameters,omitempty"`
	Signature  string        `json:"signature,omitempty"`
}
