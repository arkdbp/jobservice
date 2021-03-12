package lib

import (
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const (
	undefined = iota
	started
	success
	manualStop
	jobError
)

// Job intermediate job object
type Job struct {
	jobID     string
	path      string
	args      []string
	envs      []string
	directory string
	exitCode  int
	status    int
	output    *JobBuffer
	error     *JobBuffer
	process   *os.Process
	lock      sync.Locker
}

// JobID getter to get jobID
func (c *Job) JobID() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.jobID
}

// SetJobID will allow to set jobID
func (c *Job) SetJobID(jobID string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.jobID = jobID
}

// Process getter to get the command underlined process
func (c *Job) Process() *os.Process {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.process
}

// SetProcess will allow to set process
func (c *Job) SetProcess(process *os.Process) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.process = process
}

// ExitCode getter to get exit code
func (c *Job) ExitCode() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.exitCode
}

// SetExitCode will allow to set exit code
func (c *Job) SetExitCode(exitCode int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.exitCode = exitCode
}

// SetStatus will allow to set status
func (c *Job) SetStatus(status int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.status = status
}

// Status getter to get the status
func (c *Job) Status() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.status
}

// Output getter to get the output
func (c *Job) Output() *JobBuffer {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.output
}

// Error getter to get error
func (c *Job) Error() *JobBuffer {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.error
}

func (c *Job) getCommand() *exec.Cmd {
	c.lock.Lock()
	defer c.lock.Unlock()
	return &exec.Cmd{
		Path:        c.path,
		Args:        c.args,
		Dir:         c.directory,
		Env:         c.envs,
		Stdout:      c.output,
		Stderr:      c.error,
		SysProcAttr: &syscall.SysProcAttr{Setpgid: true},
	}
}
