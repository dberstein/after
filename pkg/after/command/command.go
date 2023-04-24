package command

import (
	"log"
	"os"
	"os/exec"
)

// Command represents a delayed execution command
type Command struct {
	cmd      *exec.Cmd
	pid      int
	exitCode int
}

// NewCommand returns new *command.DelayedCommand for `cmd` and its `args`
func NewCommand(cmd string, args ...string) *Command {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return &Command{cmd: c}
}

// Execute the command
func (c *Command) Execute() {
	if err := c.cmd.Start(); err != nil {
		log.Fatalf("%v", err)
	}
	c.pid = c.cmd.Process.Pid
	if err := c.cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			c.exitCode = exitErr.ExitCode()
		}
	}
}

// Pid returns Task's process PID or nil if not executed yet
func (c *Command) Pid() *int {
	if c.pid == 0 {
		// Process has not been executed
		return nil
	}
	return &c.pid
}

// ExitCode returns Task's process exit status or nil if not executed yet
func (c *Command) ExitCode() *int {
	if c.Pid() == nil {
		// Process has not been executed
		return nil
	}
	return &c.exitCode
}
