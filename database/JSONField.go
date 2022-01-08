package database

import (
	"bytes"
	"database/sql/driver"
	"errors"
)

type JSONField []byte

func (j JSONField) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}
func (j *JSONField) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}
func (m JSONField) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}
func (m *JSONField) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("null point exception")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// func (m *JSONField) Unmarshal(v interface{}) error {
// 	if m == nil {
// 		return errors.New("null point exception")
// 	}
//
// 	json.Unmarshal(m, v)
//
// 	// *m = append((*m)[0:0], data...)
// 	return nil
// }

func (j JSONField) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}
func (j JSONField) Equals(j1 JSONField) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}
