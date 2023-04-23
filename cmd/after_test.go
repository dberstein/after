package main

import (
	"github.com/dberstein/after/pkg/after"
	"github.com/dberstein/after/pkg/err"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func Test_isDebugValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"1"}, true},
		{"space1", args{"1 "}, true},
		{"1space", args{" 1"}, true},
		{"empty", args{""}, false},
		{"space", args{" "}, false},
		{"spaces", args{"  "}, false},
		{"0", args{"0"}, false},
		{"space0", args{" 0"}, false},
		{"0space", args{"0 "}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDebugValue(tt.args.value); got != tt.want {
				t.Errorf("isDebugValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommand(t *testing.T) {
	type args struct {
		cmdArgs []string
	}

	datePath, err := exec.LookPath("date")
	if err != nil {
		t.Fatal("ERROR: cannot find 'date' in PATH")
	}

	wantCmdDate := exec.Command(datePath, "+%FT%T")
	wantCmdDate.Stdin = os.Stdin
	wantCmdDate.Stdout = os.Stdout
	wantCmdDate.Stderr = os.Stderr

	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{"dateCmd", args{[]string{"date", "+%FT%T"}}, wantCmdDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCommand(tt.args.cmdArgs[0], tt.args.cmdArgs[1:]...); got.String() != tt.want.String() {
				t.Errorf("getCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want *err.Err
	}{
		{name: "empty", args: args{[]string{"after"}}, want: after.ErrMissingDurations},
		{name: "justDuration", args: args{[]string{"after", "1s"}}, want: after.ErrMissingCommand},
		{name: "withCmd", args: args{[]string{"after", "1s", "cmd"}}, want: nil},
		{name: "withCmdArgs", args: args{[]string{"after", "1s", "cmd", "args"}}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateArgs(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDurations(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name          string
		args          args
		wantDurations map[time.Duration]bool
		wantErr       *err.Err
	}{
		//{"invalid", args{"invalid"}, map[time.Duration]bool{}, nil},
		{"empty", args{}, nil, after.ErrMissingDurations},
		{"single", args{"*/30s"}, map[time.Duration]bool{time.Second * 0: true, time.Second * 30: true}, nil},
		{"invalid,1s", args{"invalid,1s"}, map[time.Duration]bool{time.Second * 1: true}, nil},
		{"1s,invalid", args{"1s,invalid"}, map[time.Duration]bool{time.Second * 1: true}, nil},
		{"1s", args{"1s"}, map[time.Duration]bool{time.Second * 1: true}, nil},
		{"*/20s", args{"*/20s"}, map[time.Duration]bool{time.Second * 0: true, time.Second * 20: true, time.Second * 40: true}, nil},
		{"1s,*/20s", args{"1s,*/20s"}, map[time.Duration]bool{time.Second * 0: true, time.Second * 1: true, time.Second * 20: true, time.Second * 40: true}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDurations, gotErr := getDurations(tt.args.spec)
			if !reflect.DeepEqual(gotDurations, tt.wantDurations) {
				t.Errorf("produceDurations() gotDurations = %v, want %v", gotDurations, tt.wantDurations)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("produceDurations() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_exitCode(t *testing.T) {
	type args struct {
		array []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"zeroes", args{[]int{0, 0, 0}}, 0},
		{"ones", args{[]int{1, 1, 1}}, 1},
		{"mix1", args{[]int{0, 1, 2}}, 1},
		{"mix2", args{[]int{2, 1, 0}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.args.array); got != tt.want {
				t.Errorf("exitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initDebug(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"empty", "", false},
		{"0", "0", false},
		{" 0", " 0", false},
		{"0 ", "0 ", false},
		{"1", "1", true},
		{" 1", " 1", true},
		{"1 ", "1 ", true},
		{"a", "a", true},
		{" a", " a", true},
		{"a ", "a ", true},
	}
	for _, tt := range tests {
		t.Setenv(after.DebugEnvironmentVariable, tt.value)
		t.Run(tt.name, func(t *testing.T) {
			if got := initDebug(); got != tt.want {
				t.Errorf("initDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}
