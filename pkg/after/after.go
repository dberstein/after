package after

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// InitCheck does initial check of arguments
func InitCheck(args []string) []string {
	if len(args) < 2 {
		ErrorExit(ErrMissingDurationCode, ErrMissingDurationMessage)
	}

	if len(args) < 3 {
		ErrorExit(ErrMissingCommandCode, ErrMissingCommandMessage)
	}

	return args
}

// Parse args for sleep duration considering maxDuration and returns parsed duration and remainder args
func Parse(args []string, maxDuration time.Duration) (time.Duration, []string) {
	d, err := time.ParseDuration(args[1])
	if err != nil {
		ErrorExit(ErrWrongDurationCode, err.Error())
	}

	if d > maxDuration {
		ErrorExit(ErrMaxDurationCode, fmt.Sprintf(ErrMaxDurationMessage, d, maxDuration))
	}

	return d, args[2:]
}

// Exec executes command in `args` and returns exit code
func Exec(args []string) (exitCode int) {
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("%v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return exitCode
}
