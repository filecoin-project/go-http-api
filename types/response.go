package types

type APIResponse struct {
	Errors []string `json:"errors,omitEmpty"`
}
