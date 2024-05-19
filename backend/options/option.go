package options

import "encoding/json"

type Option[T any] struct {
	Value T
	Valid bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{Value: value, Valid: true}
}

func None[T any]() Option[T] {
	return Option[T]{Valid: false}
}

func (o Option[T]) GetVal() (T, bool) {
	return o.Value, o.Valid
}

func (o Option[T]) ValueOrZero() T {
	if o.Valid {
		return o.Value
	}
	var t T
	return t
}

func (o Option[T]) IsPresent() bool {
	return o.Valid
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return json.Marshal(o.Value)
	}
	return []byte("null"), nil
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = value
	o.Valid = true
	return nil
}
