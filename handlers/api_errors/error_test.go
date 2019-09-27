package api_errors_test

import (
	"encoding/json"
	. "github.com/carbonfive/go-filecoin-rest-api/handlers/api_errors"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalErrors(t *testing.T) {
	msgs := []string{"a", "b"}

	testErr := types.APIResponse{Errors: msgs}

	testErrJSON, _ := json.Marshal(testErr)

	res := MarshalErrors(msgs)
	// len of {"errors":["a","b"]} is 20
	assert.Len(t, res, 20)
	assert.Equal(t, MarshalErrors(msgs), testErrJSON)
}
