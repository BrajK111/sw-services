package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB stores the free-form additionalDetails blob DIGIT-UI reads/writes as opaque JSON.
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("jsonb: cannot scan non-[]byte value")
	}
	return json.Unmarshal(bytes, j)
}
