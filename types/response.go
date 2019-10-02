package types

type APIErrorResponse struct {
	Errors []string `json:"errors,omitEmpty"`
}
