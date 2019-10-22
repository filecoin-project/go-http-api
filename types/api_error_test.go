package types_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-http-api/types"
)

func TestMarshalErrorStrings(t *testing.T) {
	msgs := []string{"a", "b"}
	testErr := types.APIErrorResponse{Errors: msgs}
	testErrJSON, _ := json.Marshal(testErr)
	assert.Equal(t, types.MarshalErrorStrings("a", "b"), testErrJSON)
}

func TestMarshalError(t *testing.T) {
	err := errors.New("boom")
	testErrJSON := []byte("{\"errors\":[\"boom\"]}")
	assert.Equal(t, testErrJSON, types.MarshalError(err))
}
