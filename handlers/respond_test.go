package handlers_test

import (
	"errors"
	"github.com/filecoin-project/go-http-api/types"
	"math"
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

	t.Run("error from callback responds with serialized errors & server error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		err := errors.New("boom")
		handlers.Respond(rw, nil, err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Equal(t, types.MarshalError(err), rw.Body.Bytes())
	})

	t.Run("error from marshaling responds with serialized errors and server error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		handlers.Respond(rw, math.Inf(1), nil)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func TestRespondBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("boom")
	expBody := types.MarshalError(err)

	handlers.RespondBadRequest(w, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expBody, w.Body.Bytes())
}
