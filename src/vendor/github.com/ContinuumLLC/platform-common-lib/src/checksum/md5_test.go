package checksum

import (
	"io"
	"strings"
	"testing"
)

func Test_md5Impl_Calculate(t *testing.T) {
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
			want:    "90792de52961c34118f976ebe4af3a75",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := md5Impl{}
			got, err := c.Calculate(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("md5Impl.Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("md5Impl.Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_md5Impl_Validate(t *testing.T) {
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
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792de52961c34118f976ebe4af3a75"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "validate Success Upper case",
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792DE52961C34118f976EBE4AF3A75"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "validate Failed",
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792de52961c34118f976ebe4af3a34"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "validate Trimming Success",
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792DE52961C34118f976EBE4AF3A75\n"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "validate Trimming Failed",
			args:    args{reader: strings.NewReader("Tests"), checksum: "90792DE52961C34118f976EBE4AF3A76\n"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "validate Trimming",
			args:    args{reader: strings.NewReader("Tests"), checksum: "\t \n 90792DE52961C34118f976EBE4AF3A75\n\t \t"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := md5Impl{}
			got, err := c.Validate(tt.args.reader, tt.args.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("md5Impl.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("md5Impl.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
