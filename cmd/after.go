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
	afterErr "github.com/dberstein/after/pkg/err"
)

var (
	isDebug bool
	Version string
)

func init() {
	initDebug()
}

func initDebug() bool {
	// enable `isDebug` if environment variable `after.DebugEnvironmentVariable` is debug value
	if debugValue, ok := os.LookupEnv(after.DebugEnvironmentVariable); ok && isDebugValue(debugValue) {
		isDebug = true
	}
	return isDebug
}

// validateArgs validates slice of arguments
func validateArgs(args []string) *afterErr.Err {
	switch len(args) {
	case 2:
		return afterErr.New(after.ErrMissingCommandCode, after.ErrMissingCommand.Error())
	case 1:
		return afterErr.New(after.ErrMissingDurationsCode, after.ErrMissingDurations.Error())
	default:
		return nil
	}
}

// isDebugValue return whether `value` is non-empty nor `after.DebugDisableValue`
func isDebugValue(value string) bool {
	v := strings.TrimSpace(value)
	return len(v) > 0 && v != after.DebugDisableValue
}

// getDurations returns `map[time.Duration]bool` of durations to execute based on `spec`
func getDurations(spec string) (map[time.Duration]bool, *afterErr.Err) {
	durations := after.ProduceDurations(spec)
	if len(durations) == 0 {
		return nil, after.ErrMissingDurations
	}
	return durations, nil
}

// getCommand returns `*exec.Cmd` for execution of command with arguments `cmd_args`
func getCommand(command string, args ...string) *exec.Cmd {
	// create command and wire inputs & outputs ...
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// exitCode return average of `codes`
func exitCode(codes []int) int {
	sum := 0
	for _, v := range codes {
		sum += v
	}
	return int(sum / len(codes))
}

// executeDurations executes `cmdArgs` at each `durations`
func executeDurations(durations map[time.Duration]bool, cmdArgs []string) {
	var (
		wg    sync.WaitGroup
		codes []int
	)

	if isDebug {
		cmd := getCommand(cmdArgs[0], cmdArgs[1:]...)
		_, err := fmt.Fprintf(os.Stderr, "#%s $ %s\n#%v\n", time.Now().Format(time.RFC3339Nano), cmd.String(), durations)
		if err != nil {
			panic(err)
		}
	}

	// Launch in its own goroutine each duration and wait for all to finish ...
	wg.Add(len(durations))
	for d := range durations {
		go func(d time.Duration, isDebug bool) {
			defer wg.Done()

			// sleep for requested duration before proceeding ...
			time.Sleep(d)

			cmd := getCommand(cmdArgs[0], cmdArgs[1:]...)
			if err := cmd.Start(); err != nil {
				log.Fatalf("%v", err)
			}

			if isDebug {
				_, err := fmt.Fprintf(os.Stderr, ">>%s|pid: %d|cmd: %s\n", time.Now().Format(time.RFC3339Nano), cmd.Process.Pid, cmd.String())
				if err != nil {
					panic(err)
				}
			}
			codes = append(codes, 0)
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					codes[len(codes)-1] = exitErr.ExitCode()
				}
			}
			if isDebug {
				_, err := fmt.Fprintf(os.Stderr, "<<%s|pid: %d|code: %d\n", time.Now().Format(time.RFC3339Nano), cmd.Process.Pid, codes[len(codes)-1])
				if err != nil {
					panic(err)
				}
			}
		}(d, isDebug)
	}

	// wait only if we have durations, exit with average of exit codes ...
	if len(durations) > 0 {
		wg.Wait()
		os.Exit(exitCode(codes))
	}
}

func main() {
	args := os.Args
	if len(args) > 1 && (args[1] == "-v" || args[1] == "--version") {
		fmt.Println(Version)
		return
	}

	err := validateArgs(args)
	if err != nil {
		err.Print().Exit()
	}
	durations, err := getDurations(args[1])
	if err != nil {
		err.Print().Exit()
	}
	executeDurations(durations, args[2:])
}
