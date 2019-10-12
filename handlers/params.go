package handlers

import (
	"errors"
	"net/http"
	"strings"
)

func RequireParams(r *http.Request, params ...string) error {
	errs := []string{}
	for _, param := range params {
		if valStr := r.FormValue(param); valStr == "" {
			errs = append(errs, param)
		}
	}
	if len(errs) > 0 {
		errs = append(errs, " required parameters are missing")
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}
