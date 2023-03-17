package after

import "time"

// MaxDuration is maximum duration allowed
const MaxDuration = time.Minute

// ErrMissingDurationMessage is error message for missing duration
const ErrMissingDurationMessage = "missing duration"

// ErrMissingCommandMessage is error message for missing command
const ErrMissingCommandMessage = "missing command"

// ErrMaxDurationMessage is error message for duration greater than MaxDuration
const ErrMaxDurationMessage = "duration '%s' cannot be greater than '%s'"

// ErrMissingDurationCode is exit code for missing duration
const ErrMissingDurationCode = 1

// ErrMissingCommandCode is exit code for missing command
const ErrMissingCommandCode = 2

// ErrMaxDurationCode is exit code for duration greater than MaxDuration
const ErrMaxDurationCode = 3

// ErrWrongDurationCode is exit code for invalid duration
const ErrWrongDurationCode = 4
