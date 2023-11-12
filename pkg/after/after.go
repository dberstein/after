package after

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dberstein/after/pkg/err"
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
const MinDuration = time.Millisecond

const (
	ErrMissingDurationsCode = 1 << iota
	ErrMissingCommandCode   = 1 << iota
	ErrMinMaxDuration       = 1 << iota
)

// ErrMissingDurations is error message when there are no durations
var (
	ErrMissingDurations = err.New(ErrMissingDurationsCode, "missing valid duration(s)")
	ErrMissingCommand   = err.New(ErrMissingCommandCode, "missing command")
)

func ErrMinDuration(spec string, got time.Duration) *err.Err {
	return err.Convert(ErrMinMaxDuration, fmt.Errorf("duration '%s' < '%s' (%s)", got, MinDuration, spec))
}

func ErrMaxDuration(spec string, got time.Duration) *err.Err {
	return err.Convert(ErrMinMaxDuration, fmt.Errorf("duration '%s' > '%s' (%s)", got, MaxDuration, spec))
}

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

		d, errParse := time.ParseDuration(p)
		if errParse != nil {
			_, errr := fmt.Fprintf(os.Stderr, "ERROR: %s\n", errParse.Error())
			if errr != nil {
				panic(errr)
			}
			continue
		}
		if d < MinDuration {
			ErrMinDuration(p, d).Print()
			continue
		}
		if d > MaxDuration {
			ErrMaxDuration(p, d).Print()
			continue
		}

		// if repeating duration, enable for durationDelta zero ...
		durationDelta := time.Duration(0)
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
