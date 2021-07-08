package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/matsuyoshi30/gitsu/internal/models"
)

var (
	ErrNotInsideWorktree = errors.New("not inside git worktree")
)

// SetConfig sets the user config via the 'git config' command with scope --global or --local
func SetConfig(user *models.User, scope models.Scope) error {
	err := gitConfigCommand("user.name", user.Name, scope)
	if err != nil {
		return fmt.Errorf("failed to set user.name option via git: %w", err)
	}

	err = gitConfigCommand("user.email", user.Email, scope)
	if err != nil {
		return fmt.Errorf("failed to set user.email option via git: %w", err)
	}

	err = gitGPGKeyIDCommand(user.GpgKeyID, scope)
	if err != nil {
		return fmt.Errorf("failed to set / unset user.signingkey option via git: %w", err)
	}

	return nil
}

// IsInsideWorktree returns if the current working directory is inside a git worktree
func IsInsideWorktree(scope models.Scope) error {
	// If the user sets the user config globally, we don't have to be inside a git worktree
	if scope == models.Global {
		return nil
	}

	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return err
	}

	if strings.Trim(string(out), "\n\r\t") == "true" {
		return nil
	}

	return ErrNotInsideWorktree
}

// gitConfigCommand executes a 'git config --global/--local <option> <value>' command
func gitConfigCommand(option, value string, scope models.Scope) error {
	out, err := exec.Command("git", "config", scope.Arg(), option, value).Output()
	if err != nil {
		return fmt.Errorf("%s: %w", out, err)
	}
	return nil
}

// gitGPGKeyIDCommand executes a 'git config --global/--local (--unset) user.signingkey <key ID>' command
func gitGPGKeyIDCommand(gpgKeyID string, scope models.Scope) error {
	var cmd *exec.Cmd
	if gpgKeyID == "" {
		cmd = exec.Command("git", "config", scope.Arg(), "--unset", "user.signingkey")
	} else {
		cmd = exec.Command("git", "config", scope.Arg(), "user.signingkey", gpgKeyID)
	}

	out, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 5 {
		return fmt.Errorf("%s: %w", out, err)
	}
	return nil
}
