package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

func usage() {
	format := `Usage:
  gitsu [flags]

Flags:
  --global              Set user as global.
  --gpg                 Prompt for a GPG key ID.

Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
`
	fmt.Fprintln(os.Stderr, format)
}

var (
	isGlobal  = flag.Bool("global", false, "Set user as global")
	setGpgKey = flag.Bool("gpg", false, "Prompt for a GPG key ID")

	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	os.Exit(run())
}

const (
	sel = "Select git user"
	add = "Add new git user"
	del = "Delete git user"
)

func run() int {
	action := promptui.Select{
		Label: "Select action",
		Items: []string{sel, add, del},
	}

	_, actionType, err := action.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to select action: %v\n", err)
		return 1
	}

	switch actionType {
	case sel:
		if err := selectUser(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select user: %v\n", err)
			return 1
		}
	case add:
		if err := addUser(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add user: %v\n", err)
			return 1
		}
	case del:
		if err := deleteUser(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete user: %v\n", err)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "Unexpected action type\n")
		return 1
	}

	return 0
}

func selectUser() error {
	users, err := ListUser()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users")
		return nil
	}

	user := promptui.Select{
		Label: "Select git user",
		Items: UsersToString(users),
	}
	selectedUserIndex, _, err := user.Run()
	if err != nil {
		return err
	}

	option := "--local"
	if *isGlobal {
		option = "--global"
	}

	cmdName := exec.Command("git", "config", option, "user.name", users[selectedUserIndex].Name)
	if err := cmdName.Run(); err != nil {
		return err
	}
	cmdMail := exec.Command("git", "config", option, "user.email", users[selectedUserIndex].Email)
	if err := cmdMail.Run(); err != nil {
		return err
	}
	if users[selectedUserIndex].GpgKeyID != "" {
		cmdGpgKey := exec.Command("git", "config", option, "user.signingkey", users[selectedUserIndex].GpgKeyID)
		if err := cmdGpgKey.Run(); err != nil {
			return err
		}
	}

	return nil
}

func addUser() error {
	name := promptui.Prompt{
		Label: "Input git user name",
	}
	resultName, err := name.Run()
	if err != nil {
		return err
	}

	email := promptui.Prompt{
		Label:    "Input git email address",
		Validate: ValidateEmail,
	}
	resultEmail, err := email.Run()
	if err != nil {
		return err
	}
	var resultKeyID string
	if *setGpgKey {
		keyIdPrompt := promptui.Prompt{
			Label: "Input GPG key ID",
		}
		resultKeyID, err = keyIdPrompt.Run()
		if err != nil {
			return err
		}
	}

	if err := CreateUser(resultName, resultEmail, resultKeyID); err != nil {
		return err
	}

	return nil
}

func deleteUser() error {
	users, err := ListUser()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users")
		return nil
	}

	user := promptui.Select{
		Label: "Select git user",
		Items: UsersToString(users),
	}
	selectedUserIndex, _, err := user.Run()
	if err != nil {
		return err
	}

	if err := RemoveUser(selectedUserIndex, users); err != nil {
		return err
	}

	return nil
}
