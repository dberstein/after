package task

import (
	"github.com/dberstein/after/pkg/after/command"
	"time"
)

// Task represents a single task execution
type Task struct {
	cmd *command.Command
}

// NewTask returns a new *Task
func NewTask(cmd string, args ...string) *Task {
	return &Task{cmd: command.NewCommand(cmd, args...)}
}

// Execute after `delay` the Task's command
func (t *Task) Execute(delay time.Duration) int {
	time.Sleep(delay)
	t.cmd.Execute()
	return *t.cmd.ExitCode()
}

// Cmd returns *command.Command
func (t *Task) Cmd() command.Command {
	return *t.cmd
}

// Pid returns execution PID. Zero means task not yet executed
func (t *Task) Pid() int {
	return *t.cmd.Pid()
}

// ExitCode returns exit code of Task's process PID (see command.Pid)
func (t *Task) ExitCode() int {
	return *t.cmd.ExitCode()
}
