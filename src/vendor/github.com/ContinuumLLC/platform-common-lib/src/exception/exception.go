package exception

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

// Exception interface is the enterprise error interface.
// Use it especially for handling nesting and wrapping of errors.
// It implements the in-built 'error' interface.
type Exception interface {
	error
	// GetInnerError returns the wrapped error object.
	GetInnerError() error
	// GetErrorCode returns the error code of this error.
	GetErrorCode() string
	// Get Data Variables stored with this error
	GetData() map[string]interface{}
}

// New function creates an exception with the specified error code and inner-error
func New(errCode string, innerErr error) Exception {
	return newException(errCode, innerErr, nil)
}

func newException(errCode string, innerErr error, dataVars map[string]interface{}) *exceptionImpl {
	return &exceptionImpl{
		errorCode: errCode,
		inner:     innerErr,
		callers:   callers(),
		data:      dataVars,
	}
}

func NewWithMap(errCode string, innerErr error, dataVars map[string]interface{}) Exception {
	return newException(errCode, innerErr, dataVars)
}

type exceptionImpl struct {
	errorCode string
	inner     error
	callers   *stack
	data      map[string]interface{}
}

func (exc *exceptionImpl) GetInnerError() error {
	return exc.inner
}

func (exc *exceptionImpl) GetErrorCode() string {
	return exc.errorCode
}

func (exc *exceptionImpl) Error() string {
	return exc.errorCode
}

func (exc *exceptionImpl) GetData() map[string]interface{} {
	return exc.data
}

// This number needs to be 4 because the exception package "New" method itself makes a call to an private method.
// Before the creation of NewWithMap method, this value was 4.
const cIgnoreInitialCallersCount = 4

func callers() *stack {
	const depth = 32
	pcs := make(stack, depth)
	n := runtime.Callers(cIgnoreInitialCallersCount, pcs[:])
	pcs = pcs[0:n]
	return &pcs
}

func (exc *exceptionImpl) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintln(s, exc.Error())
		verbStr := "%v"
		if s.Flag('+') {
			verbStr = "%+v"
		}
		if exc.data != nil && len(exc.data) > 0 {
			fmt.Fprintf(s, "\tdata-"+verbStr, exc.data)
			fmt.Fprintln(s)
		}
		exc.callers.Format(s, verb)
		if exc.inner != nil {
			fmt.Fprintf(s, verbStr, exc.inner)
			fmt.Fprintln(s)
		}
	case 's':
		fmt.Fprintln(s, exc.Error())
	}
}

type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		if len(*s) > 0 {
			s.formatStack(st, (*s)[0])
		}
		switch {
		case st.Flag('+'):
			for _, pc := range (*s)[1:] {
				s.formatStack(st, pc)
			}
		}
	}
}

func (s *stack) formatStack(st fmt.State, pc uintptr) {
	f := errors.Frame(pc)
	fmt.Fprintf(st, "\t\t")
	fmt.Fprintf(st, "%+v", f)
	fmt.Fprintln(st)
}
