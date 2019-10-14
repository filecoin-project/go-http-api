package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-http-api/types"
)

func TestMarshalErrors(t *testing.T) {
	msgs := []string{"a", "b"}

	testErr := types.APIErrorResponse{Errors: msgs}

	testErrJSON, _ := json.Marshal(testErr)

	assert.Equal(t, types.MarshalErrors(msgs), testErrJSON)
}
