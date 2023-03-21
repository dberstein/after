package main

import (
	"fmt"
	"github.com/dberstein/after/pkg/after"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
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

	var (
		wg        sync.WaitGroup
		exitCodes int
	)

	durations := after.ProduceDurations(os.Args[1])
	if len(durations) == 0 {
		_, errDurations := fmt.Fprintf(os.Stderr, "ERROR: "+after.ErrMissingDurationsMessage+"\n")
		if errDurations != nil {
			panic(errDurations)
		}
		os.Exit(after.ErrMissingDurationCode)
	}

	// debug info if environment variable `after.DebugEnvironmentVariable` is non-empty nor `after.DebugDisableValue` ...
	if debug, ok := os.LookupEnv(after.DebugEnvironmentVariable); ok && len(strings.TrimSpace(debug)) > 0 && strings.TrimSpace(debug) != after.DebugDisableValue {
		_, err := fmt.Fprintf(os.Stderr, "%v\n: %s\n", durations, strings.Join(os.Args[2:], " "))
		if err != nil {
			panic(err)
		}
	}

	// Launch in its own goroutine each duration and wait for all to finish ...
	wg.Add(len(durations))
	for d := range durations {
		go func(d time.Duration) {
			// sleep for requested duration before proceeding ...
			time.Sleep(d)

			// create command and wire inputs & outputs ...
			cmd := exec.Command("sh", "-c", strings.Join(os.Args[2:], " "))
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				log.Fatalf("%v", err)
			}

			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCodes += exitErr.ExitCode()
				}
			}

			defer wg.Done()
		}(d)
	}

	// wait only if we have durations, exit with average of exit codes ...
	if len(durations) > 0 {
		wg.Wait()
		os.Exit(exitCodes / len(durations))
	}
}
