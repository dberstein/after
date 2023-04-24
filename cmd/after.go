package main

import (
	"fmt"
	"github.com/dberstein/after/pkg/after/tasks"
	"github.com/dberstein/after/pkg/color"
	"os"
	"path"
	"strings"
)

var (
	Version string
)

func version() {
	fmt.Println(Version)
}

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "%sUsage%s: %s%s ( -h | --help | <duration(s)> <command> [args...] )%s\n", color.Green, color.Reset, color.Cyan, prog, color.Reset)
}

func parse(args ...string) int {
	switch len(args) {
	case 1:
		usage(path.Base(args[0]))
		return 1
	case 2:
		usage(path.Base(args[0]))
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			return 0
		} else {
			return 1
		}
	default:
		return -1
	}
}

func main() {
	switch parse(os.Args...) {
	case 0:
		os.Exit(0)
	case 1:
		fmt.Fprintf(os.Stderr, "%s%s%s\n", color.Red, "not enough parameters", color.Reset)
		os.Exit(1)
	default:
		cmd := os.Args[2]
		args := os.Args[3:]
		ts := tasks.NewTasks(cmd, args...)
		for _, spec := range strings.Split(os.Args[1], ",") {
			ts.Add(spec)
		}
		os.Exit(ts.Execute())
	}
}
