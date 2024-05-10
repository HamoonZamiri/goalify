package responses

type APIError struct {
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
}

func NewAPIError(message string, errors map[string]string) *APIError {
	return &APIError{
		Message: message,
		Errors:  errors,
	}
}
