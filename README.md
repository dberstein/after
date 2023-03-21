# after
Cronjob utility to target sub-minute times

![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/dberstein/after/go.yml?branch=main) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dberstein/after) ![GitHub](https://img.shields.io/github/license/dberstein/after) ![GitHub last commit](https://img.shields.io/github/last-commit/dberstein/after)

## Usage

    $ crontab -l
    * * * * * after <duration(s)> <command [args]>

- `<duration(s)>` must be comma-separated list of durations, as understood by <a href="https://pkg.go.dev/time#ParseDuration">Go</a>, and shorter than one minute. For durations over one minute use regular <a href="https://en.wikipedia.org/wiki/Cron#Overview">Cron</a> spec.
  - Durations can also be of the repeating form `*/<duration>` which will repeat every `<duration>` within the same minute. Example: `*/20s` will run on seconds `0`, `20`, and `40` of the minute.
  - Durations can be combined in a comma separated list, like: `5s,*/20s500ms,15s`
- `<command [args]>` must be command and optional arguments to execute.
  - `<command [args]>` is only executed once per concurrent durations, meaning `*/15,*/30` will NOT run command twice at seconds `0` and `30` although both expressions coincide in those seconds.

## Standard in, out, err

- `<command>` receives `stdin`, `stdout` and `stderr` from `after`. 

## Exit codes

Exit code will be exit code of required `<command>`, except in these cases:

- `1` missing duration(s)
- `2` missing command

## Build

### Requirements

- GNU make
- Golang

Binaries are built as `./build/<os>-<arch>/after`:

    $ make build

## Install

Builds and installs `./build/<os>-<arch>/after` as `/usr/bin/after`

    $ sudo make install [-e INSTALL_DIR=/usr/bin]

### Uninstall

    $ sudo make uninstall [-e INSTALL_DIR=/usr/bin]


## Examples

    # Every 15m run and log twice "date" at seconds 20 and 45 of the minute
    */15 * * * * after 20s,45S date >> date1.log

    # Every 15m run and every 5 seconds and at second 33 with 500 milliseconds log "date"
    */15 * * * * after '*/5,33s500ms' date >> date2.log
