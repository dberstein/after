package main

import (
	"os"
	"time"

	"github.com/dberstein/after/pkg/after"
)

var args []string

func init() {
	args = after.InitCheck(os.Args)
}

func main() {
	d, args, err := after.Parse(args, after.MaxDuration)
	time.Sleep(d)
	exitCode := after.Exec(args)
	os.Exit(exitCode)
}
