package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type SerialisableStringArray []string

func (s SerialisableStringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *SerialisableStringArray) Scan(src interface{}) (err error) {
	var stringArray []string
	err = json.Unmarshal([]byte(src.(string)), &stringArray)
	if err != nil {
		return
	}
	*s = stringArray
	return nil
}
