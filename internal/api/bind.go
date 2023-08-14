package api

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

type BindError struct {
	err error
}

func NewBindError(err error) *BindError {
	return &BindError{err: err}
}

func (e *BindError) Error() string {
	return e.err.Error()
}

func (e *BindError) Unwrap() error {
	return e.err
}

func Bind(c echo.Context, v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		panic("v must be pointer")
	}

	if err := c.Bind(v); err != nil {
		return NewBindError(err)
	}
	return nil
}
