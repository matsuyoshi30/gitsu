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
	mod = "Modify git user"
)

func run() int {
	action := promptui.Select{
		Label:  "Select action",
		Items:  []string{sel, add, del},
		Stdout: &bellSkipper{},
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
	case mod:
		if err := modifyUser(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to modify user: %v\n", err)
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

	userPrompt := promptui.Select{
		Label: "Select git user",
		Items: UsersToString(users),
	}
	selectedUserIndex, _, err := userPrompt.Run()
	if err != nil {
		return err
	}

	user := users[selectedUserIndex]

	scopeOpt := "--local"
	if *isGlobal {
		scopeOpt = "--global"
	}

	cmdName := exec.Command("git", "config", scopeOpt, "user.name", user.Name)
	if err := cmdName.Run(); err != nil {
		return err
	}
	cmdMail := exec.Command("git", "config", scopeOpt, "user.email", user.Email)
	if err := cmdMail.Run(); err != nil {
		return err
	}

	cmdGpgKey := exec.Command("git", "config", scopeOpt, "user.signingkey", user.GpgKeyID)
	if users[selectedUserIndex].GpgKeyID == "" {
		cmdGpgKey = exec.Command("git", "config", scopeOpt, "--unset", "user.signingkey")
	}
	if err := cmdGpgKey.Run(); err != nil {
		// git exits with code 5 when unsetting a non-existent property
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 5 {
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

func modifyUser() error {
	users, err := ListUser()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users")
		return nil
	}

	userPrompt := promptui.Select{
		Label: "Select git user",
		Items: UsersToString(users),
	}
	selectedUserIndex, _, err := userPrompt.Run()
	if err != nil {
		return err
	}

	name := promptui.Prompt{
		Label: "Input git user name, leave empty for no change",
	}
	newName, err := name.Run()
	if err != nil {
		return err
	}

	email := promptui.Prompt{
		Label:    "Input git email address, leave empty for no change",
		Validate: ValidateEmail,
	}
	newEmail, err := email.Run()
	if err != nil {
		return err
	}

	var newKeyID string
	if *setGpgKey {
		keyIdPrompt := promptui.Prompt{
			Label: "Input GPG key ID, leave empty for no change",
		}
		newKeyID, err = keyIdPrompt.Run()
		if err != nil {
			return err
		}
	}

	newUser := User{
		Name:     newName,
		Email:    newEmail,
		GpgKeyID: newKeyID,
	}

	return ModifyUser(selectedUserIndex, newUser)
}
