package gormutil

import (
	"database/sql/driver"
	"encoding/json"
)

type StringList []string

func (g StringList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *StringList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
