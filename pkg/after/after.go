package after

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// RepeatPrefix is prefix for repeating durations
const RepeatPrefix = "*/"

// DebugEnvironmentVariable is environmental variable to control debug output
const DebugEnvironmentVariable = "DEBUG"

// DebugDisableValue is special value to disable debug output
const DebugDisableValue = "0"

// MaxDuration is maximum duration allowed
const MaxDuration = time.Minute

// MinDuration is minimum duration allowed
const MinDuration = time.Duration(0)

// ErrMissingDurationsMessage is error message when there are no durations
const ErrMissingDurationsMessage = "missing valid duration(s)"

// ErrMissingCommandMessage is error message for missing command
const ErrMissingCommandMessage = "missing command"

// ErrMaxDurationMessage is error message for duration greater than MaxDuration
const ErrMaxDurationMessage = "duration '%s' cannot be greater than '%s'"

// ErrMinDurationMessage is error message for duration smaller than MinDuration
const ErrMinDurationMessage = "duration '%s' cannot be less than '%s'"

// ErrMissingDurationCode is exit code for missing duration
const ErrMissingDurationCode = 1

// ErrMissingCommandCode is exit code for missing command
const ErrMissingCommandCode = 2

// ProduceDurations produce durations from strings like "<d>", "<d>,<d>", "*/<d>" and its combinations.
func ProduceDurations(spec string) map[time.Duration]bool {
	durations := make(map[time.Duration]bool)

	// first split by commas
	for _, p := range strings.Split(spec, ",") {
		// whether this is a repeating duration or not ...
		repeatDuration := strings.HasPrefix(p, RepeatPrefix)
		if repeatDuration {
			p = strings.TrimPrefix(p, RepeatPrefix)
		}

		d, err := time.ParseDuration(p)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
			if err != nil {
				panic(err)
			}
			continue
		}
		if d < MinDuration {
			_, err = fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(ErrMinDurationMessage, d, MinDuration))
			continue
		}
		if d > MaxDuration {
			_, err = fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(ErrMaxDurationMessage, d, MaxDuration))
			continue
		}

		// if repeating duration, enable for durationDelta zero ...
		durationDelta := MinDuration
		if repeatDuration && d < MaxDuration-time.Nanosecond {
			durations[durationDelta] = true
		}

		// loop until durationDelta is too big ...
		for durationDelta < MaxDuration-time.Nanosecond-d {
			// increase durationDelta and enable it ...
			durationDelta += d
			durations[durationDelta] = true

			// ensure we won't repeat if we shouldn't ...
			if !repeatDuration {
				durationDelta += MaxDuration
			}
		}
	}

	return durations
}
