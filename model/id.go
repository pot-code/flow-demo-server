package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type UUID uint64

func (id *UUID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(*id), 10)), nil
}

func (id *UUID) UnmarshalJSON(b []byte) error {
	s := strings.Replace(string(b), "\"", "", -1)
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*id = UUID(u)
	return nil
}

// Scan assigns a value from a database driver.
// The src value will be of one of the following types:
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
//	nil - for NULL values
//
// An error should be returned if the value cannot be stored
// without loss of information.
//
// Reference types such as []byte are only valid until the next call to Scan
// and should not be retained. Their underlying memory is owned by the driver.
// If retention is necessary, copy their values before the next call to Scan.
func (id *UUID) Scan(src any) error {
	switch t := src.(type) {
	case int64:
		*id = UUID(t)
	case []byte:
		v, err := strconv.ParseUint(string(t), 10, 64)
		if err != nil {
			return fmt.Errorf("parse uuid: %w", err)
		}
		*id = UUID(v)
	}
	return nil
}

// Value returns a driver Value.
// Value must not panic.
func (id *UUID) Value() (driver.Value, error) {
	return int64(*id), nil
}