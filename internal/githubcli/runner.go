package githubcli

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Runner is the small command boundary used by the GitHub adapter. Tests use a
// fake runner; production uses the authenticated gh CLI already required by the
// project-administration workflow.
type Runner interface {
	Run(ctx context.Context, args ...string) ([]byte, error)
}

type inputRunner interface {
	RunInput(ctx context.Context, input []byte, args ...string) ([]byte, error)
}

// ExecRunner invokes gh without a shell.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, args ...string) ([]byte, error) {
	return runGH(ctx, nil, args...)
}

// RunInput invokes gh with exact bytes on standard input. It is used for API
// endpoints whose JSON request body cannot be represented safely as flags.
func (ExecRunner) RunInput(ctx context.Context, input []byte, args ...string) ([]byte, error) {
	return runGH(ctx, input, args...)
}

func runGH(ctx context.Context, input []byte, args ...string) ([]byte, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, errors.New("GitHub CLI (gh) is not installed; install gh and authenticate it before using GitHub commands")
	}
	command := exec.CommandContext(ctx, "gh", args...)
	if input != nil {
		command.Stdin = strings.NewReader(string(input))
	}
	output, err := command.Output()
	if err == nil {
		return output, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		message := strings.TrimSpace(string(exitError.Stderr))
		if message != "" {
			return nil, fmt.Errorf("gh %s: %s", strings.Join(args, " "), message)
		}
	}
	return nil, fmt.Errorf("gh %s: %w", strings.Join(args, " "), err)
}
