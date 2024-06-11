package responses

type ServerResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func New[T any](data T, message string) *ServerResponse[T] {
	return &ServerResponse[T]{
		Data:    data,
		Message: message,
	}
}
