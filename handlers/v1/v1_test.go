package v1_test

import (
	"fmt"
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCallbacks_BuildHandlers(t *testing.T) {
	// check that all possible callbacks can be set
	// check that they will get a default handler
	cb := v1.Callbacks{}

	cbt := reflect.TypeOf(cb)
	numCallbacks := cbt.NumField()
	for i := 0; i < numCallbacks; i++ {
		fieldName := cbt.Field(i).Name
		test1Name := fmt.Sprintf("%s will get a default handler", fieldName)

		t.Run(test1Name, func(t *testing.T) {
			hlers := *cb.BuildHandlers()
			assert.NotNil(t, hlers[fieldName])
			assert.Equal(t, &v1.DefaultHandler{}, hlers[fieldName])
		})
	}
}
