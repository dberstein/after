package main

import (
	"fmt"
	"github.com/dberstein/after/pkg/after/tasks"
	"os"
	"strings"
)

var (
	Version string
)

func version() {
	fmt.Println(Version)
}

func init() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s\n", "not enough parameters")
		os.Exit(1)
	}
}
func main() {
	cmd := os.Args[2]
	args := os.Args[3:]
	ts := tasks.NewTasks(cmd, args...)
	for _, spec := range strings.Split(os.Args[1], ",") {
		ts.Add(spec)
	}
	os.Exit(ts.Execute())
}
