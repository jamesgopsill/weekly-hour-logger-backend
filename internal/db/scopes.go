package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	SYS_ADMIN_SCOPE = "sysadmin"
	ADMIN_SCOPE     = "admin"
	USER_SCOPE      = "user"
)

// https://stackoverflow.com/questions/41375563/unsupported-scan-storing-driver-value-type-uint8-into-type-string

type Scopes []string

func (s Scopes) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *Scopes) Scan(src interface{}) (err error) {
	var scopes []string
	err = json.Unmarshal([]byte(src.(string)), &scopes)
	if err != nil {
		return
	}
	*s = scopes
	return nil
}
