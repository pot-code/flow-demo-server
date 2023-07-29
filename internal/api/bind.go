package api

import (
	"fmt"

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
	if err := c.Bind(v); err != nil {
		return NewBindError(fmt.Errorf("数据解析失败，请检查输入"))
	}
	return nil
}
