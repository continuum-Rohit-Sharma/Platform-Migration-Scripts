package tasking

import (
	"encoding/json"
	"fmt"
)

// Regularity type is used for regularity definition
type Regularity int

// These constants describe regularity of the task
const (
	_ Regularity = iota
	RunNow
	OneTime
	Recurrent
)

// UnmarshalJSON is used to Unmarshal Regularity from JSON
func (regularity *Regularity) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}
	return regularity.Parse(stringValue)
}

// Parse is used to Parse string Regularity to Regularity type
func (regularity *Regularity) Parse(s string) error {
	switch s {
	case "":
		*regularity = 0
	case "RunNow":
		*regularity = RunNow
	case "OneTime":
		*regularity = OneTime
	case "Recurrent":
		*regularity = Recurrent
	default:
		return fmt.Errorf("incorrect regularity: %s", s)
	}
	return nil
}

// MarshalJSON custom marshal method for field Regularity
func (regularity Regularity) MarshalJSON() ([]byte, error) {
	switch regularity {
	case 0:
		return json.Marshal("")
	case RunNow:
		return json.Marshal("RunNow")
	case OneTime:
		return json.Marshal("OneTime")
	case Recurrent:
		return json.Marshal("Recurrent")
	default:
		return []byte{}, fmt.Errorf("incorrect task regularity: %v", regularity)
	}
}
