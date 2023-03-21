package after

import (
	"reflect"
	"testing"
	"time"
)

func TestProduceDurations(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name string
		args args
		want map[time.Duration]bool
	}{
		{
			name: "1s",
			args: args{"1s"},
			want: map[time.Duration]bool{time.Second: true},
		},
		{
			name: "-1s,30s,59s",
			args: args{"-1s,30s,59s"},
			want: map[time.Duration]bool{time.Second * 30: true, time.Second * 59: true},
		},
		{
			name: "1invalid,30s,61s",
			args: args{"1invalid,30s,61s"},
			want: map[time.Duration]bool{time.Second * 30: true},
		},
		{
			name: "*/20s",
			args: args{"*/20s"},
			want: map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 40: true},
		},
		{
			name: "*/20s,20s",
			args: args{"*/20s,20s"},
			want: map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 40: true},
		},
		{
			name: "*/20s,21s",
			args: args{"*/20s,21s"},
			want: map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 21: true, time.Second * 40: true},
		},
		{
			name: "21s,*/20s",
			args: args{"21s,*/20s"},
			want: map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 21: true, time.Second * 40: true},
		},
		{
			name: "21s,*/20s,21s",
			args: args{"21s,*/20s,21s"},
			want: map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 21: true, time.Second * 40: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProduceDurations(tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProduceDurations() = %v, want %v", got, tt.want)
			}
		})
	}
}
