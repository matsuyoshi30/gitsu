package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func SetGitConfig(user User, scope Scope, users []User, selectedUserIndex int) error {
	cmdName := exec.Command("git", "config", scope.String(), "user.name", user.Name)
	if err := cmdName.Run(); err != nil {
		return err
	}

	cmdMail := exec.Command("git", "config", scope.String(), "user.email", user.Email)
	if err := cmdMail.Run(); err != nil {
		return err
	}

	cmdGpgKey := exec.Command("git", "config", scope.String(), "user.signingkey", user.GpgKeyID)
	if users[selectedUserIndex].GpgKeyID == "" {
		cmdGpgKey = exec.Command("git", "config", scope.String(), "--unset", "user.signingkey")
	}
	if err := cmdGpgKey.Run(); err != nil {
		// git exits with code 5 when unsetting a non-existent property
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 5 {
			return err
		}
	}

	return nil
}

func IsUnderGitDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(cwd, ".git")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("you need to be in a git project directory to use gitsu")
	} else if err != nil {
		return err
	}

	cmdRevParse := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmdRevParse.Run(); err != nil {
		return err
	}

	return nil
}
