package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB wraps a map for PostgreSQL jsonb columns.
type JSONB map[string]any

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid jsonb value: %T", value)
	}
	if len(bytes) == 0 {
		*j = JSONB{}
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// StringList wraps []string for PostgreSQL jsonb columns.
type StringList []string

func (s *StringList) Scan(value any) error {
	if value == nil {
		*s = StringList{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid jsonb list value: %T", value)
	}
	if len(bytes) == 0 {
		*s = StringList{}
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}
