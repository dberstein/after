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

### Quoting

If command or arguments to be executed by `after` require quoting, please use `printf %q` to escape them, likeso:


    $ after \*/5s sh -c $(printf %q 'sleep 2 && echo "escap'\''ed: $(date)"')
    escap'ed: Sat Mar 25 07:43:25 IDT 2023
    escap'ed: Sat Mar 25 07:43:30 IDT 2023
    escap'ed: Sat Mar 25 07:43:35 IDT 2023
    ...

See [debug](#debug) for more diagnostic output, for example:

    $ DEBUG=1 after \*/20s sh -c $(printf %q 'sleep 22 && echo "escap'\''ed: $(date)"')
    [cmd: /bin/sh -c sh -c sleep\ 22\ \&\&\ echo\ \"escap\'ed:\ \$\(date\)\"]
    @map[0s:true 20s:true 40s:true]
    >>pid: 65247 | cmd: /bin/sh -c sh -c sleep\ 22\ \&\&\ echo\ \"escap\'ed:\ \$\(date\)\"
    >>pid: 65257 | cmd: /bin/sh -c sh -c sleep\ 22\ \&\&\ echo\ \"escap\'ed:\ \$\(date\)\"
    escap'ed: Sat Mar 25 07:51:05 IDT 2023
    <<pid: 65247 | code: 0
    >>pid: 65269 | cmd: /bin/sh -c sh -c sleep\ 22\ \&\&\ echo\ \"escap\'ed:\ \$\(date\)\"
    escap'ed: Sat Mar 25 07:51:25 IDT 2023
    <<pid: 65257 | code: 0
    escap'ed: Sat Mar 25 07:51:45 IDT 2023
    <<pid: 65269 | code: 0

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

Builds binary `./build/<OS>-<ARCH>/after` and installs it as binary `/usr/bin/after`. You can control installation directory with environmental variable or make parameter `INSTALL_DIR` wich has default value `/usr/bin`.

    $ sudo make install [-e INSTALL_DIR=/usr/bin]

### Uninstall

Uninstalls binary of `make install`. If install used custom `INSTALL_DIR`, same value must be used for `make uninstall` to remove correct binary.

    $ sudo make uninstall [-e INSTALL_DIR=/usr/bin]

### Debug

If environmental variable `DEBUG` has a non-empty and different from `0` value, debug information is sent to `stderr`. Information includes full command line being executed, the schedule of execution and each execution's exit code.

For example compare:

    $ after \*/20s date +%T
    08:33:21
    08:33:41
    08:34:01

With, where schedule of executions (0s, 20s, 40s) and each execution and its exit code are displayed in stderr:

    $ DEBUG=1 after \*/20s date +%T
    [cmd: /bin/sh -c date +%T]
    @map[0s:true 20s:true 40s:true]
    >>pid: 69963 | cmd: /bin/sh -c date +%T
    08:26:35
    <<pid: 69963 | code: 0
    >>pid: 69972 | cmd: /bin/sh -c date +%T
    08:26:55
    <<pid: 69972 | code: 0
    >>pid: 69989 | cmd: /bin/sh -c date +%T
    08:27:15
    <<pid: 69989 | code: 0

## Examples

### Fixed seconds

Every 15m run and log twice "date" at seconds 20 and 45 of the minute:

    */15 * * * * after 20s,45S date >> date1.log

### Repeating seconds

Every 15m run and every 5 seconds and at second 33 with 500 milliseconds log "date". Note that `*/...` as duration for `after` must be quoted or escaped:

    */15 * * * * after '*/5,33s500ms' date >> date2.log
