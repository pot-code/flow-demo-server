package api

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func DecodeFromRequestBody(receiver any, body io.ReadCloser) error {
	if reflect.TypeOf(receiver).Kind() != reflect.Ptr {
		panic("receive must be a pointer")
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return NewDecoderError(fmt.Errorf("read request body: %w", err))
	}

	if err := json.Unmarshal(data, receiver); err != nil {
		return NewDecoderError(fmt.Errorf("unmarshal request body: %w", err))
	}

	return nil
}

type DecoderError struct {
	err error
}

func NewDecoderError(err error) *DecoderError {
	return &DecoderError{err: err}
}

func (s DecoderError) Error() string {
	return s.err.Error()
}

func (s DecoderError) Unwrap() error {
	return s.err
}
