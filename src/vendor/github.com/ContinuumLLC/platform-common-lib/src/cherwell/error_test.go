package cherwell

import (
	"reflect"
	"testing"
)

func TestNewErrorSet(t *testing.T) {

	tests := []struct {
		name string
		es   []error
		want *Errors
	}{
		{name: "New set of errors",
			es:   []error{&GeneralFailure{Message: "GeneralFailure"}, &RecordNotFound{Message: "RecordNotFound"}},
			want: NewErrorSet(&GeneralFailure{Message: "GeneralFailure"}, &RecordNotFound{Message: "RecordNotFound"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorSet(tt.es...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrorSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorsError(t *testing.T) {

	tests := []struct {
		name   string
		Errors []error
		want   string
	}{
		{name: "Check Error() for Errors",
			Errors: []error{&GeneralFailure{Message: "GeneralFailure"}, &RecordNotFound{Message: "RecordNotFound"}},
			want:   "GeneralFailure\nRecordNotFound",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Errors{
				Errors: tt.Errors,
			}
			if got := es.Error(); got != tt.want {
				t.Errorf("Errors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
