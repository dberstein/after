package main

import (
	"github.com/dberstein/after/pkg/after"
	"os"
	"time"
)

var args []string

func init() {
	args = after.InitCheck(os.Args)
}

func main() {
	d, args := after.Parse(args, after.MaxDuration)
	time.Sleep(d)
	exitCode := after.Exec(args)
	os.Exit(exitCode)
}
