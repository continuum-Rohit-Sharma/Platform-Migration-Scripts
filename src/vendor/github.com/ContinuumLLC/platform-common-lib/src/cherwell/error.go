package cherwell

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BusObNotValidError constant for checking Cherwell responses on valid data
const BusObNotValidError = "BusObNotValid"

// RecordNotFoundError constant for checking Cherwell responses on missing objects
const RecordNotFoundError = "RECORDNOTFOUND"

// GeneralFailureError constant for checking Cherwell responses on different error
const GeneralFailureError = "GENERALFAILURE"

// ErrorData holds common error info in response
type ErrorData struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	HasError     bool   `json:"hasError"`
}

// GetErrorObject gets error object based on error response
func (r *ErrorData) GetErrorObject() error {
	switch r.ErrorCode {
	case BusObNotValidError:
		return &BusObNotValid{Message: r.ErrorMessage}
	case RecordNotFoundError:
		return &RecordNotFound{Message: r.ErrorMessage}
	case GeneralFailureError:
		return &GeneralFailure{Message: r.ErrorMessage}
	default:
		return &CherwellError{Code: r.ErrorCode, Message: r.ErrorMessage}
	}

	return nil
}

// CherwellError describes a cherwell api error
type CherwellError struct {
	error
	Code    string
	Message string
}

func (ge *CherwellError) Error() string {
	return fmt.Sprintf(ge.Message)
}

// Errors is a set of errors
type Errors struct {
	Errors []error
}

// NewErrorSet creates a new Errors instance
func NewErrorSet(es ...error) *Errors {
	return &Errors{Errors: es}
}

func (es *Errors) Error() string {
	errs := make([]string, len(es.Errors))
	for i, err := range es.Errors {
		errs[i] = err.Error()
	}

	return strings.Join(errs, "\n")
}

// Add appends an error or set of errors to existing set
func (es *Errors) Add(e ...error) {
	es.Errors = append(es.Errors, e...)
}

// IsEmpty checks if set of errors is not empty
func (es *Errors) IsEmpty() bool{
	return len(es.Errors) == 0
}

type BusObNotValid struct {
	error
	Message string
}

func (ge BusObNotValid) Error() string {
	return fmt.Sprintf(ge.Message)
}

type RecordNotFound struct {
	error
	Message string
}

func (ge RecordNotFound) Error() string {
	return fmt.Sprintf(ge.Message)
}

type GeneralFailure struct {
	error
	Message string
}

func (ge GeneralFailure) Error() string {
	return fmt.Sprintf(ge.Message)
}

type InvalidFilterOperator struct {
	Message string
}

func (ifo InvalidFilterOperator) Error() string {
	return fmt.Sprintf(ifo.Message)
}

func containsNotFound(message string) bool {
	return strings.Contains(message, "not found")
}

// errorFromResponse builds error based on error response
func errorFromResponse(body string) error {
	errBody := ErrorData{}
	if err := json.Unmarshal([]byte(body), &errBody); err != nil {
		return &GeneralFailure{
			Message: body,
		}
	}

	return errBody.GetErrorObject()
}
