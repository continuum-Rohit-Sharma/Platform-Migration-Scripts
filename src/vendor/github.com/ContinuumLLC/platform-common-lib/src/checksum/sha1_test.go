package checksum

import (
	"io"
	"strings"
	"testing"
)

func Test_sha1Impl_Calculate(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "calculate Success",
			args:    args{strings.NewReader("Tests")},
			want:    "39fdec1194d94212b871a28b2aa04a73cd40fce1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := sha1Impl{}
			got, err := c.Calculate(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("sha1Impl.Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sha1Impl.Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sha1Impl_Validate(t *testing.T) {
	type args struct {
		reader   io.Reader
		checksum string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "validate Success",
			args:    args{reader: strings.NewReader("Tests"), checksum: "39fdec1194d94212b871a28b2aa04a73cd40fce1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "validate Success upper case",
			args:    args{reader: strings.NewReader("Tests"), checksum: "39FDEC1194D94212b871A28B2aa04a73cd40fce1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "validate Failed",
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792de52961c34118f976ebe4af3a34"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := sha1Impl{}
			got, err := c.Validate(tt.args.reader, tt.args.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("sha1Impl.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sha1Impl.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
