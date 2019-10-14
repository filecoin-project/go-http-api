package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-http-api/handlers"
)

func TestRespond(t *testing.T) {
	t.Run("sets status code and serializes the result", func(t *testing.T) {
		type TestResult struct {
			Data string
		}

		rw := httptest.NewRecorder()
		handlers.Respond(rw, &TestResult{"abcd"}, nil)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, `{"Data":"abcd"}`, rw.Body.String())
	})

	t.Run("responds with serialized errors", func(t *testing.T) {
		rw := httptest.NewRecorder()
		handlers.Respond(rw, nil, errors.New("boom"))

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Equal(t, `{"errors":["boom"]}`, rw.Body.String())
	})
}