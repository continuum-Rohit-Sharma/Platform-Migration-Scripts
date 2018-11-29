package tasking

import (
	"encoding/json"
	"fmt"
)

// TaskState type is used for states definition
type TaskState int

// These constants describe states of the task
const (
	_ TaskState = iota
	TaskStateActive
	TaskStateInactive
	TaskStateDisabled
)

// UnmarshalJSON is used to Unmarshal Task State from JSON
func (state *TaskState) UnmarshalJSON(byteResult []byte) error {
	var stringValue string
	if err := json.Unmarshal(byteResult, &stringValue); err != nil {
		return err
	}
	return state.Parse(stringValue)
}

// Parse is used to Parse string TaskState to TaskState type
func (state *TaskState) Parse(s string) error {
	switch s {
	case "":
		*state = TaskState(0)
	case "Active":
		*state = TaskStateActive
	case "Inactive":
		*state = TaskStateInactive
	case "Disabled":
		*state = TaskStateDisabled
	default:
		return fmt.Errorf("incorrect state: %s", s)
	}
	return nil
}

// MarshalJSON custom marshal method for field State
func (state TaskState) MarshalJSON() ([]byte, error) {
	switch state {
	case 0:
		return json.Marshal("")
	case TaskStateActive:
		return json.Marshal("Active")
	case TaskStateInactive:
		return json.Marshal("Inactive")
	case TaskStateDisabled:
		return json.Marshal("Disabled")
	default:
		return []byte{}, fmt.Errorf("incorrect task state: %v", state)
	}
}
