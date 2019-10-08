package v1_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/carbonfive/go-filecoin-rest-api/handlers/v1"
	"github.com/stretchr/testify/assert"
)

func TestRespond(t *testing.T) {
	t.Run("sets status code and serializes the result", func(t *testing.T) {
		type TestResult struct {
			Data string
		}

		rw := httptest.NewRecorder()
		v1.Respond(rw, &TestResult{"abcd"}, nil)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, `{"Data":"abcd"}`, rw.Body.String())
	})

	t.Run("responds with serialized errors", func(t *testing.T) {
		rw := httptest.NewRecorder()
		v1.Respond(rw, nil, errors.New("boom"))

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, `{"errors":["boom"]}`, rw.Body.String())
	})
}
