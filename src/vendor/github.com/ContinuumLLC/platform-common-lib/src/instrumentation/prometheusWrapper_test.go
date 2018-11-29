package instrumentation

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

type testLogger struct{}

func (testLogger) Infof(format string, v ...interface{})  {}
func (testLogger) Errorf(format string, v ...interface{}) {}

func TestStartListening(t *testing.T) {
	type args struct {
		enableListener bool
		logger         logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nilLogger",
			args: args{
				enableListener: false,
				logger:         nil,
			},
			wantErr: true,
		},
		{
			name: "StartServerFalse",
			args: args{
				enableListener: false,
				logger:         testLogger{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartListening(tt.args.enableListener, tt.args.logger); (err != nil) != tt.wantErr {
				t.Errorf("StartListening() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterNewMetricsCounter(t *testing.T) {
	type args struct {
		name string
		help string
	}
	tests := []struct {
		name string
		args args
		want prometheus.Counter
	}{
		{
			name: "default",
			args: args{
				name: "testCounter",
				help: "testCounter Help",
			},
			want: prometheus.NewCounter(prometheus.CounterOpts{Name: "testCounter", Help: "testCounter Help"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegisterNewMetricsCounter(tt.args.name, tt.args.help); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterNewMetricsCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterNewGaugeCounter(t *testing.T) {
	type args struct {
		name string
		help string
	}
	tests := []struct {
		name string
		args args
		want prometheus.Gauge
	}{
		{
			name: "default",
			args: args{
				name: "testGaugeCounter",
				help: "testCounter Help",
			},
			want: prometheus.NewGauge(prometheus.GaugeOpts{Name: "testGaugeCounter", Help: "testCounter Help"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegisterNewGaugeCounter(tt.args.name, tt.args.help); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterNewGaugeCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucket(t *testing.T) {
	type args struct {
		start float64
		width float64
		count int
	}
	tests := []struct {
		name string
		args args
		want Bucket
	}{
		{
			name: "default",
			args: args{
				start: 10.0,
				width: 20.0,
				count: 5,
			},
			want: Bucket{
				Start: 10.0,
				Width: 20.0,
				Count: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucket(tt.args.start, tt.args.width, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterNewHistrogram(t *testing.T) {
	type args struct {
		name   string
		help   string
		bucket Bucket
	}
	tests := []struct {
		name string
		args args
		want prometheus.Histogram
	}{
		{
			name: "default",
			args: args{
				name:   "testHistCounter",
				help:   "testCounter Help",
				bucket: NewBucket(1.0, 2.0, 5),
			},
			want: prometheus.NewHistogram(prometheus.HistogramOpts{Name: "testHistCounter", Help: "testCounter Help", Buckets: prometheus.LinearBuckets(1.0, 2.0, 5)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegisterNewHistrogram(tt.args.name, tt.args.help, tt.args.bucket); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterNewHistrogram() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestRegisterNewSummary(t *testing.T) {
// 	type args struct {
// 		name       string
// 		help       string
// 		objectives map[float64]float64
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want prometheus.Summary
// 	}{
// 		{
// 			name: "default",
// 			args: args{
// 				name:       "testCounter",
// 				help:       "testCounter Help",
// 				objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
// 			},
// 			want: prometheus.NewSummary(prometheus.SummaryOpts{Name: "testCounter", Help: "testCounter Help", Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}}),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := RegisterNewSummary(tt.args.name, tt.args.help, tt.args.objectives); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RegisterNewSummary() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
