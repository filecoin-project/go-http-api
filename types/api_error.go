package types

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
)

type APIErrorResponse struct {
	Errors []string `json:"errors,omitempty"`
}

func IrrecoverableError() APIErrorResponse {
	return APIErrorResponse{Errors: []string{"irrecoverable error"}}
}

func MarshalErrors(errlist []string) []byte {
	errs := APIErrorResponse{Errors: errlist}
	res, err := json.Marshal(errs)
	if err != nil {
		res, _ = json.Marshal(IrrecoverableError())
	}
	return res
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
