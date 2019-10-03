package api_errors

import (
	"encoding/json"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

func IrrecoverableError() types.APIErrorResponse {
	return types.APIErrorResponse{Errors: []string{"irrecoverable error"}}
}

func MarshalErrors(errlist []string) []byte {
	errs := types.APIErrorResponse{Errors: errlist}
	res, err := json.Marshal(errs)
	if err != nil {
		res, _ = json.Marshal(IrrecoverableError())
	}
	return res
}
