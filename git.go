package main

import "os/exec"

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
