package main

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type User struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
	Alias string `json:"Alias`
}

func UsersToString(users []User) []string {
	us := make([]string, len(users))
	for i, user := range users {
		us[i] = user.Name + " <" + user.Email + ">" + " Alias: " + user.Alias
	}

	return us
}

func ListUser() ([]User, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	if !IsExist(configPath) {
		return nil, errors.New("No user")
	}

	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	return config.Users, nil
}

func CreateUser(name, email, alias string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	if alias == "" {
		alias = name
	}

	if !IsExist(configPath) {
		if err := CreateConfig(Config{Users: []User{User{Name: name, Email: email, Alias: alias}}}); err != nil {
			return err
		}
		return nil
	}

	config, err := ReadConfig()
	if err != nil {
		return err
	}

	if config.Users == nil {
		config.Users = []User{User{Name: name, Email: email, Alias: alias}}
	} else {
		config.Users = append(config.Users, User{Name: name, Email: email, Alias: alias})
	}

	if err := CreateConfig(config); err != nil {
		return err
	}

	return nil
}

func RemoveUser(idx int, users []User) error {
	newUsers := make([]User, len(users)-1)
	if idx+1 == len(users) {
		newUsers = users[:idx]
	} else {
		newUsers = append(users[:idx], users[idx+1:]...)
	}

	if err := CreateConfig(Config{Users: newUsers}); err != nil {
		return err
	}

	return nil
}

func ModifyUser(idx int, users []User, modifiedUser User) error {
	config, err := ReadConfig()
	
	users = config.Users

	if err != nil {
		return err
	}

	if modifiedUser.Name != "" {
		users[idx].Name = modifiedUser.Name
	}

	if modifiedUser.Email != "" {
		users[idx].Email = modifiedUser.Email
	}


	if modifiedUser.Alias != "" {
		users[idx].Alias = modifiedUser.Alias
	}


	if err := CreateConfig(Config{Users: users}); err != nil {
		return err
	}
	return nil

}

func ValidateEmail(email string) error {
	if !govalidator.IsExistingEmail(email) {
		return errors.New("Invalid email address")
	}

	return nil
}
