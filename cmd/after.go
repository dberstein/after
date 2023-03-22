package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/dberstein/after/pkg/after"
)

func main() {
	if len(os.Args) < 2 {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", after.ErrMissingDurationsMessage)
		if err != nil {
			panic(err)
		}
		os.Exit(after.ErrMissingDurationCode)
	}

	if len(os.Args) < 3 {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", after.ErrMissingCommandMessage)
		if err != nil {
			panic(err)
		}
		os.Exit(after.ErrMissingCommandCode)
	}

	executeDurations(produceDurations(os.Args[1]), os.Args[2:])
}

// produceDurations return `map[time.Duration]bool` of durations to execute based on `spec`
func produceDurations(spec string) (durations map[time.Duration]bool) {
	durations = after.ProduceDurations(spec)
	if len(durations) == 0 {
		_, errDurations := fmt.Fprintf(os.Stderr, "ERROR: "+after.ErrMissingDurationsMessage+"\n")
		if errDurations != nil {
			panic(errDurations)
		}
		os.Exit(after.ErrMissingDurationCode)
	}

	return
}

// executeDurations executes `cmd_args` at each
func executeDurations(durations map[time.Duration]bool, cmd_args []string) {
	var (
		wg              sync.WaitGroup
		exitCodes, code int
		isDebug         bool
	)

	// debug info if environment variable `after.DebugEnvironmentVariable` is non-empty nor `after.DebugDisableValue` ...
	if debug, ok := os.LookupEnv(after.DebugEnvironmentVariable); ok && len(strings.TrimSpace(debug)) > 0 && strings.TrimSpace(debug) != after.DebugDisableValue {
		isDebug = true
	}

	if isDebug {
		cmd := getCommand(cmd_args)
		_, err := fmt.Fprintf(os.Stderr, "[%s]@%v\n", cmd.String(), durations)
		if err != nil {
			panic(err)
		}
	}

	// Launch in its own goroutine each duration and wait for all to finish ...
	wg.Add(len(durations))

	for d := range durations {
		go func(d time.Duration, isDebug bool) {
			defer wg.Done()
			cmd := getCommand(cmd_args)

			// sleep for requested duration before proceeding ...
			time.Sleep(d)

			if err := cmd.Start(); err != nil {
				log.Fatalf("%v", err)
			}

			code = 0
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					code = exitErr.ExitCode()
				}
			}
			exitCodes += code

			if isDebug {
				_, err := fmt.Fprintf(os.Stderr, "Execution [%s] (code: %d)\n", cmd.String(), code)
				if err != nil {
					panic(err)
				}
			}
		}(d, isDebug)
	}

	// wait only if we have durations, exit with average of exit codes ...
	if len(durations) > 0 {
		wg.Wait()
		os.Exit(exitCodes / len(durations))
	}
}

// getCommand returns `*exec.Cmd` for execution of command with arguments `cmd_args`
func getCommand(cmd_args []string) *exec.Cmd {
	// create command and wire inputs & outputs --...
	cmd := exec.Command("sh", "-c", strings.Join(cmd_args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
