package timeutil

import (
	"reflect"
	"testing"
	"time"
)

func TestToHourLong(t *testing.T) {
	fromTm1, _ := time.Parse(time.RFC3339, "2018-02-05T11:24:53.489110939Z")
	toTm1, _ := time.Parse(time.RFC3339, "2018-02-05T12:24:53.489110939Z")

	fromTm2, _ := time.Parse(time.RFC3339, "2018-02-05T11:24:53.489110939Z")
	toTm2, _ := time.Parse(time.RFC3339, "2018-02-05T20:44:53.489110939Z")

	fromTm3, _ := time.Parse(time.RFC3339, "2018-02-05T23:24:53.489110939Z")
	toTm3, _ := time.Parse(time.RFC3339, "2018-02-06T01:24:53.489110939Z")

	fromTm4, _ := time.Parse(time.RFC3339, "2018-02-04T23:24:53.489110939Z")
	toTm4, _ := time.Parse(time.RFC3339, "2018-02-06T05:24:53.489110939Z")

	type args struct {
		fromTime time.Time
		toTime   time.Time
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test1",
			args: args{fromTime: fromTm1, toTime: toTm1},
			want: []int{20180205},
		},
		{
			name: "Test2",
			args: args{fromTime: fromTm2, toTime: toTm2},
			want: []int{20180205},
		},
		{
			name: "Test3",
			args: args{fromTime: fromTm3, toTime: toTm3},
			want: []int{20180205, 20180206},
		},
		{
			name: "Test4",
			args: args{fromTime: fromTm4, toTime: toTm4},
			want: []int{20180204, 20180205, 20180206},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToHourLong(tt.args.fromTime, tt.args.toTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToHourLong() = %v, want %v", got, tt.want)
			}
		})
	}
}
