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
  gitsu [flags] <alias>

Flags:
  --global              Set user as global.

  alias                 Change username/email/alias for given alias
Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
`
	fmt.Fprintln(os.Stderr, format)
}

var (
	isGlobal = flag.Bool("global", false, "Set user as global")

	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) > 1{
		alias := os.Args[1]
		if alias != "" {
			modifyUser(alias)
		}
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
	action := promptui.Select{
		Label: "Select action",
		Items: []string{sel, add, del, mod},
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
		if err := modifyUser(""); err != nil {
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

	alias := promptui.Prompt{
		Label: "Input git user alias, leave empty to set it for username",
	}

	resultAlias, err := alias.Run()
    if err != nil {
		return err
	}

	if err := CreateUser(resultName, resultEmail, resultAlias); err != nil {
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

func modifyUser(userAlias string) error {
	users, err := ListUser()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users")
		return nil
	}


	var selectedUserIndex = -1
	for i, user := range users {
		if user.Alias == userAlias{
			selectedUserIndex = i
		}
	}

	if selectedUserIndex < 0 {
		user := promptui.Select{
			Label: "Select git user",
			Items: UsersToString(users),
		}
		selectedUserIndex, _, err = user.Run()
		if err != nil {
			return err
		}
	}

	name := promptui.Prompt{
		Label: "Input git user name",
	}
	resultName, err := name.Run()
	if err != nil {
		return err
	}

	email := promptui.Prompt{
		Label:    "Input git email address",
	}

	resultEmail, err := email.Run()
	if err != nil {
		return err
	}

	if resultEmail != "" {
		if err := ValidateEmail(resultEmail); err != nil {
			return err
		}
	}

	alias := promptui.Prompt{
		Label: "Input git user alias, leave empty to set it without change",
	}

	resultAlias, err := alias.Run()
    if err != nil {
		return err
	}

	modifiedUser := User{Name: resultName, Email: resultEmail, Alias: resultAlias}
	if err := ModifyUser(selectedUserIndex, users, modifiedUser); err != nil {
		return err
	}

	return nil
}
