package api_errors

import (
	"encoding/json"
	"github.com/carbonfive/go-filecoin-rest-api/types"
)

func IrrecoverableError() types.APIResponse {
	return types.APIResponse{Errors: []string{"irrecoverable error"}}
}

func MarshalErrors(errlist []string) []byte {
	errs := types.APIResponse{Errors: errlist}
	res, err := json.Marshal(errs)
	if err != nil {
		res, _ = json.Marshal(IrrecoverableError())
	}
	return res
}
