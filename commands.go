package main

import (
	"fmt"
	"os/exec"

	"github.com/manifoldco/promptui"
)

// command is a function type of a command performed by our service.
type command func() error

// commands is a map that contains all supported commands and their description
// as keys.
var commands = map[string]command{
	"Select git user":  selectUser,
	"Add new git user": addUser,
	"Delete git user":  deleteUser,
}

// listCommands returns a slice of strings that contains descriptions for all
// supported commands.
func listCommands() (slice []string) {
	for k := range commands {
		slice = append(slice, k)
	}
	return
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

	if err := CreateUser(resultName, resultEmail); err != nil {
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
