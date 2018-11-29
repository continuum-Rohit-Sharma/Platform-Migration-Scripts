package autoupdate

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Value - Implementation of valuer for database/sql
func (s Status) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type such as string
	return string(s), nil
}

// Scan - Implement the database/sql scanner interface
func (s *Status) Scan(value interface{}) error {
	// if value is nil, Failure
	if value == nil {
		// set the value of the pointer s to Failure
		*s = Failure
		return nil
	}
	if str, err := driver.String.ConvertValue(value); err == nil {
		// set the value of the pointer s to Status(v)
		*s = Status(fmt.Sprintf("%s", str))
		return nil

	}
	// otherwise, return an error
	return errors.New("Failed to scan Status")
}

// Value - Implementation of valuer for database/sql
func (i InstallVariables) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type such as string
	return json.Marshal(i)
}

// Scan - Implement the database/sql scanner interface
func (i *InstallVariables) Scan(value interface{}) error {
	// if value is nil, Failure
	if value == nil {
		// set the value of the pointer s to Failure
		*i = InstallVariables{}
		return nil
	}
	if v, ok := value.([]byte); ok {
		// set the value of the pointer s to Status(v)
		*i = InstallVariables{}
		return json.Unmarshal(v, &i)
	}
	// otherwise, return an error
	return errors.New("Failed to scan Install Variables")
}
