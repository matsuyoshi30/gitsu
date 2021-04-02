package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

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
	alias     = ""
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

	if len(os.Args) > 1 {
		alias = os.Args[1]
	}

	os.Exit(run())
}

const (
	sel = "Select git user"
	add = "Add new git user"
	del = "Delete git user"
	mod = "Modify git user"
)

func run() int {
	if alias != "" {
		if err := selectSingleUser(alias); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select user: %v\n", err)
			return 1
		}

		return 0
	}

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

	scope := Local
	if *isGlobal {
		scope = Global
	}

	return SetGitConfig(user, scope, users, selectedUserIndex)
}

func selectSingleUser(alias string) error {
	users, err := ListUser()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users")
		return nil
	}

	scope := Local
	if *isGlobal {
		scope = Global
	}

	for idx, user := range users {
		if user.Alias == alias {
			return SetGitConfig(user, scope, users, idx)
		}
	}

	msg := fmt.Sprintf("No user with alias [%s]", alias)
	return errors.New(msg)
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

	alias := promptui.Prompt{
		Label: "Input git user alias, leave empty for no alias",
	}
	resultAlias, err := alias.Run()
	if err != nil {
		return err
	}

	if err := CreateUser(resultName, resultEmail, resultKeyID, resultAlias); err != nil {
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
		Validate: ValidateModifiedEmail,
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

	alias := promptui.Prompt{
		Label: "Input git user alias, leave empty for no alias",
	}
	newAlias, err := alias.Run()
	if err != nil {
		return err
	}

	newUser := User{
		Name:     newName,
		Email:    newEmail,
		Alias:    newAlias,
		GpgKeyID: newKeyID,
	}

	return ModifyUser(selectedUserIndex, newUser)
}
