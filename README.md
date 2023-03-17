# after
Cronjob utility to target sub-minute times

## Usage

    $ crontab -l
    * * * * * after <duration> <command>

- `<duration>` must be any duration understood by <a href="https://pkg.go.dev/time#ParseDuration">Go</a> and shorter than one minute. For durations over one minute use regular <a href="https://en.wikipedia.org/wiki/Cron#Overview">Cron</a> spec.
- `<command>` must be command and optional arguments to execute.

## Standard in, out, err

- `<command>` receives `stdin`, `stdout` and `stderr` from `after`. 

## Exit codes

Exit code will be exit code of required `<command>`, except in these cases:

- `1` missing duration
- `2` missing command
- `3` duration is greater then one minute
- `4` duration could not be parsed

## Build

### Requirements

- GNU make
- Golang


    $ make build

Binaries are built as `./build/<os>-<arch>/after`

## Install

Builds and installs `./build/<os>-<arch>/after` as `/usr/local/bin/after`

    $ make build && sudo make install

## Examples

    # Every 15m run and log twice "date" at seconds 20 and 45 of the minute
    */15 * * * * after 20s date >> date.log
    */15 * * * * after 45s date >> date.log