// Package git provides functions which interact with the git binary to set / unset config values
package git

import (
	"fmt"
	"os/exec"

	"github.com/matsuyoshi30/gitsu/internal/models"
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

	return nil
}

// IsInsideWorktree returns if the current working directory is inside a git worktree
func IsInsideWorktree() (bool, error) {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false, err
	}

	if string(out) == "true" {
		return true, nil
	}
	return false, nil
}

// gitConfigCommand executes a 'git config --global/--local <option> <value>' command
func gitConfigCommand(option, value string, scope models.Scope) error {
	out, err := exec.Command("git", "config", scope.Arg(), option, value).Output()
	if err != nil {
		return fmt.Errorf("%s: %w", out, err)
	}
	return nil
}
