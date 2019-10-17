package types

import (
	"errors"
	"net/http"
)

// SignedMessage is a struct representing a complete signed message in binary format,
// sent to the HTTP API as Content-Type: application/octet-stream
type SignedMessage struct {
	MessageBlob []byte `json:"messageBlob,required"`
}

func (sm *SignedMessage) Bind(r *http.Request) error {
	if len(sm.MessageBlob) == 0 {
		return errors.New("messageBlob is missing")
	}
	return nil
}
