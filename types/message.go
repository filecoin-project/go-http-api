package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"
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

func (m *Message) BindRequest(r *http.Request) error {
	var err error

	m.To = r.FormValue("to")
	vstr := r.FormValue("value")
	var ok bool
	m.Value, ok = big.NewInt(0).SetString(vstr, 10)
	if !ok {
		return errors.New("failed to parse big.Int: Value")
	}
	m.GasPrice, ok = big.NewInt(0).SetString(r.FormValue("gasPrice"), 10)
	if !ok {
		return errors.New("failed to parse big.Int: GasPrice")
	}
	if m.GasLimit, err = strconv.ParseUint(r.FormValue("gasLimit"), 10, 64); err != nil {
		return err
	}
	m.Method = r.FormValue("method")
	m.Parameters = strings.Split(r.FormValue("parameters"), ",")
	return nil
}
