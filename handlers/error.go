package handlers

import "encoding/json"

type ErrStruct struct {
	error string
}

// MarshalError marshals an error returned from a filecoin callback function
func MarshalError(cberr error) string {
	body, err := json.Marshal(cberr)
	if err != nil {
		return IrrecoverableError()
	}
	return string(body[:])
}

func IrrecoverableError() string {
	msg, _ := json.Marshal(ErrStruct{error: "irrecoverable error"})
	return string(msg[:])
}
