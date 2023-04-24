package tasks

import (
	"fmt"
	"github.com/dberstein/after/pkg/after/task"
	"github.com/dberstein/after/pkg/color"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const PrefixMany = "*/"

// Tasks represents c cmd, args that should be executed at different time.Duration
type Tasks struct {
	cmd   string
	args  []string
	tasks map[time.Duration]*task.Task
}

// NewTasks returns new Tasks struct
func NewTasks(cmd string, args ...string) *Tasks {
	return &Tasks{cmd: cmd, args: args, tasks: make(map[time.Duration]*task.Task)}
}

// Add a new time spec to tasks' schedule
func (ts *Tasks) Add(spec string) *Tasks {
	if strings.HasPrefix(spec, PrefixMany) {
		return ts.addMany(strings.TrimPrefix(spec, PrefixMany))
	}
	return ts.addOne(spec)
}

func (ts *Tasks) parseDuration(spec string) *time.Duration {
	d, err := time.ParseDuration(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", color.Red, err, color.Reset)
		return nil
	}
	if d < time.Duration(0) {
		fmt.Fprintf(os.Stderr, "%s%s%s: %s%s%s\n", color.Red,
			fmt.Sprintf("duration cannot be less than %s", time.Duration(0)),
			color.Reset, color.Yellow, spec, color.Reset)
		return nil
	} else if d >= time.Minute {
		fmt.Fprintf(os.Stderr, "%s%s%s: %s%s%s\n", color.Red,
			fmt.Sprintf("duration cannot be greater than %s", time.Minute),
			color.Reset, color.Yellow, spec, color.Reset)
		return nil
	}
	return &d
}

func (ts *Tasks) addOne(spec string) *Tasks {
	if d := ts.parseDuration(spec); d != nil {
		ts.tasks[*d] = task.NewTask(ts.cmd, ts.args...)
	}
	return ts
}

func (ts *Tasks) addMany(spec string) *Tasks {
	if d := ts.parseDuration(spec); d != nil {
		ds := time.Duration(0)
		for ds < time.Minute {
			ts.tasks[ds] = task.NewTask(ts.cmd, ts.args...)
			ds += *d
		}
	}
	return ts
}

// Execute loops through all tasks and executes async all of them
func (ts *Tasks) Execute() int {
	var wg sync.WaitGroup
	wg.Add(len(ts.tasks))
	for d := range ts.tasks {
		go func(tsk *task.Task, d time.Duration) {
			defer wg.Done()
			tsk.Execute(d)
		}(ts.tasks[d], d)
	}
	wg.Wait()
	return ts.ExitCode()
}

// ExitCode is aggregated exit code average of executed tasks
func (ts *Tasks) ExitCode() int {
	sum := 0
	for d := range ts.tasks {
		t := ts.tasks[d]
		if t.Pid() > 0 { // if already executed ...
			sum += t.ExitCode()
		}
	}
	avg := sum / int(math.Max(1, float64(len(ts.tasks))))
	return avg
}
