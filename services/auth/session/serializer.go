package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type SessionSerializer interface {
	Serialize(s *Session) ([]byte, error)
	Deserialize(b []byte) (*Session, error)
}

type serializer struct{}

func NewSerializer() *serializer {
	return &serializer{}
}

func (se *serializer) Serialize(s *Session) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(s)
	if err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}
	return b.Bytes(), nil
}

func (se *serializer) Deserialize(b []byte) (*Session, error) {
	s := new(Session)
	err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(s)
	if err != nil {
		return nil, fmt.Errorf("gob decode: %w", err)
	}
	return s, nil
}
