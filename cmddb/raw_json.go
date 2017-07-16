package cmddb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONRaw is a custom type for JSONB fields in db
type JSONRaw json.RawMessage

// Value converts json to driver.Value
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)

	return driver.Value(byteArr), nil
}

// Scan allows scanning of JSON data
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		if len(*j) == 0 {
			return nil
		}
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}

	return nil
}

// MarshalJSON converts json to byte slice
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte(nil), nil
	}
	return *j, nil
}

// UnmarshalJSON converts byte slice to json
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
