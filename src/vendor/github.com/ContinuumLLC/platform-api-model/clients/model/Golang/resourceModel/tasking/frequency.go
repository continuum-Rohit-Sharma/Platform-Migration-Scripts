package tasking

import (
	"encoding/json"
	"fmt"
)

//Frequency basis for recurrent execution
type Frequency int

// These constants describe frequency of recurrent execution
const (
	_ Frequency = iota
	Hourly
	Daily
	Weekly
	Monthly
)

// UnmarshalJSON used to convert string Frequency representation to Frequency type
func (frequency *Frequency) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}
	switch stringValue {
	case "":
		*frequency = 0
	case "hourly":
		*frequency = Hourly
	case "daily":
		*frequency = Daily
	case "weekly":
		*frequency = Weekly
	case "monthly":
		*frequency = Monthly
	default:
		return fmt.Errorf("incorrect frequency: %s", stringValue)
	}

	return nil
}

// MarshalJSON custom marshal method for field Frequency
func (frequency Frequency) MarshalJSON() ([]byte, error) {
	switch frequency {
	case 0:
		return json.Marshal("")
	case Hourly:
		return json.Marshal("hourly")
	case Daily:
		return json.Marshal("daily")
	case Weekly:
		return json.Marshal("weekly")
	case Monthly:
		return json.Marshal("monthly")
	default:
		return []byte{}, fmt.Errorf("incorrect task frequency : %v", frequency)
	}
}
