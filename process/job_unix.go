//go:build !windows

package process

import (
	"os"
)

// Job is noop on unix
type Job int

var (
	// JobObject is public global JobObject, 0 value on linux
	JobObject Job
)

// CreateJobObject returns a job object.
func CreateJobObject() (pj Job, err error) {
	return pj, err
}

// NewJob is noop on unix
func NewJob() (Job, error) {
	return 0, nil
}

// Close is noop on unix
func (job Job) Close() error {
	return nil
}

// Assign is noop on unix
func (job Job) Assign(p *os.Process) error {
	return nil
}
