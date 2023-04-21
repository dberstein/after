# after
Cronjob utility to target sub-minute times

![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/dberstein/after/go.yml?branch=main) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dberstein/after) ![GitHub](https://img.shields.io/github/license/dberstein/after) ![GitHub last commit](https://img.shields.io/github/last-commit/dberstein/after)

## Usage

    $ crontab -l
    * * * * * after <duration(s)> <command [args]>

- `<duration(s)>` must be comma-separated list of durations, as understood by <a href="https://pkg.go.dev/time#ParseDuration">Go</a>, at least `1ms` and shorter than one minute (`1m`). For durations over one minute use regular <a href="https://en.wikipedia.org/wiki/Cron#Overview">Cron</a> spec.
  - Durations can also be of the repeating form `*/<duration>` which will repeat every `<duration>` within the same minute. Example: `*/20s` will run on seconds `0`, `20`, and `40` of the minute.
  - Durations can be combined in a comma separated list, like: `5s,*/20s500ms,15s`
- `<command [args]>` must be command and optional arguments to execute.
  - `<command [args]>` is only executed once per concurrent durations, meaning `*/15,*/30` will NOT run command twice at seconds `0` and `30` although both expressions coincide in those seconds.

### Quoting

If command or arguments to be executed by `after` require quoting, please use quotes and/or escape them, like so:

    $ after \*/20s sh -c "sleep 2 && echo escape\'d \$(date)"
    escape'd Sat Apr 1 07:55:29 IDT 2023
    escape'd Sat Apr 1 07:55:49 IDT 2023
    escape'd Sat Apr 1 07:56:09 IDT 2023

See [debug](#debug) for more diagnostic output, for example:

    $ DEBUG=1 after \*/20s sh -c "sleep 2 && echo escape\'d \$(date)"
    #2023-04-01T07:53:52.599851+03:00 $ /bin/sh -c sh -c "sleep 2 && echo escape\\'d \$(date)"
    #map[0s:true 20s:true 40s:true]
    >>2023-04-01T07:53:52.602151+03:00|pid: 52770|cmd: /bin/sh -c sh -c "sleep 2 && echo escape\\'d \$(date)"
    escape'd Sat Apr 1 07:53:54 IDT 2023
    <<2023-04-01T07:53:54.627701+03:00|pid: 52770|code: 0
    >>2023-04-01T07:54:12.606055+03:00|pid: 52784|cmd: /bin/sh -c sh -c "sleep 2 && echo escape\\'d \$(date)"
    escape'd Sat Apr 1 07:54:14 IDT 2023
    <<2023-04-01T07:54:14.631721+03:00|pid: 52784|code: 0
    >>2023-04-01T07:54:32.604656+03:00|pid: 52798|cmd: /bin/sh -c sh -c "sleep 2 && echo escape\\'d \$(date)"
    escape'd Sat Apr 1 07:54:34 IDT 2023
    <<2023-04-01T07:54:34.635794+03:00|pid: 52798|code: 0

## Standard in, out, err

- `<command [args]>` receives `stdin`, `stdout` and `stderr` from `after`.

## Exit codes

Exit code will be exit code of `<command [args]>`, in addition of these cases related to `after` operation:

- `1` missing duration(s). No valid duration was given either first parameter missing or none of the values given is valid (comma-separated list).
- `2` missing command. Second parameter onwards must be command to execute when duration(s) expire.

## Build

### Requirements

- GNU make
- Golang

Binaries are built as binary `./build/<OS>-<ARCH>/after`:

    $ make build

## Install

Builds binary `./build/<OS>-<ARCH>/after` and installs it as binary `/usr[/local]/bin/after`. 

    $ sudo make install

### Uninstall

Uninstalls binary of `make install`.

    $ sudo make uninstall

### Debug

If environmental variable `DEBUG` has a non-empty and different from `0` value, debug information is sent to `stderr`.

Information includes full command line being executed, the schedule of execution and each execution's process (pid), timestamp and exit code.

For example compare:

    $ after \*/20s date +%T
    07:59:45
    08:00:05
    08:00:25

With, where schedule of executions (0s, 20s, 40s) and each execution and its exit code are displayed in STDERR:

    $ DEBUG=1 after \*/20s date +%T
    #2023-04-01T07:59:45.884795+03:00 $ /bin/sh -c date +%T
    #map[0s:true 20s:true 40s:true]
    >>2023-04-01T07:59:45.886871+03:00|pid: 53094|cmd: /bin/sh -c date +%T
    07:59:45
    <<2023-04-01T07:59:45.895521+03:00|pid: 53094|code: 0
    >>2023-04-01T08:00:05.892511+03:00|pid: 53105|cmd: /bin/sh -c date +%T
    08:00:05
    <<2023-04-01T08:00:05.906056+03:00|pid: 53105|code: 0
    >>2023-04-01T08:00:25.891189+03:00|pid: 53124|cmd: /bin/sh -c date +%T
    08:00:25
    <<2023-04-01T08:00:25.905354+03:00|pid: 53124|code: 0

## Examples

### Fixed seconds

Every 15m run and log twice "date" at seconds 20 and 45 of the minute:

    */15 * * * * after 20s,45S date >> date1.log

### Repeating seconds

Every 15m run and every 5 seconds and at second 33 with 500 milliseconds log "date". Note that `*/...` as duration for `after` must be quoted or escaped:

    */15 * * * * after '*/5,33s500ms' date >> date2.log
