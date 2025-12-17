package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes external commands while honoring dry-run mode.
type Runner struct {
	DryRun bool
}

// Cmd runs the provided command and returns captured stdout/stderr.
// When DryRun is true the command is not executed.
func (r Runner) Cmd(name string, args ...string) (string, string, error) {
	cmdline := strings.Join(append([]string{name}, args...), " ")
	if r.DryRun {
		return "", "", nil
	}

	cmd := exec.Command(name, args...) // #nosec G204 -- arguments provided by trusted CLI flags.
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("%s failed: %w", cmdline, err)
	}

	return stdout.String(), stderr.String(), nil
}
