// Package launchd asks launchd about the user's jobs.
package launchd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// LinkLabel is the link daemon's job.
const LinkLabel = "com.remote-agents.link"

// jobNotFound is launchctl's exit status for a job launchd does not have.
const jobNotFound = 113

type Launchctl struct {
	Path string
	UID  int
}

// Loaded reports whether launchd has the job in the user's GUI domain. It
// has a job from login on, before it starts the job's process.
func (l Launchctl) Loaded(ctx context.Context, label string) (bool, error) {
	out, err := exec.CommandContext(ctx, l.Path, "print", fmt.Sprintf("gui/%d/%s", l.UID, label)).CombinedOutput()
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		return true, nil
	case errors.As(err, &exitErr) && exitErr.ExitCode() == jobNotFound:
		return false, nil
	}
	return false, fmt.Errorf("asking launchd for %s: %w: %s", label, err, out)
}
