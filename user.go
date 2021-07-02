package main

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type User struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Alias    string `json:"Alias"`
	GpgKeyID string `json:"GpgKeyId"`
}

func UsersToString(users []User) []string {
	us := make([]string, len(users))
	for i, user := range users {
		alias := ""
		if user.Alias != "" {
			alias = fmt.Sprintf("[%s] ", user.Alias)
		}

		us[i] = fmt.Sprintf("%s%s <%s>", alias, user.Name, user.Email)
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

func CreateUser(name, email, gpgKeyID, alias string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	user := User{
		Name:     name,
		Email:    email,
		GpgKeyID: gpgKeyID,
		Alias:    alias,
	}

	if !IsExist(configPath) {
		if err := CreateConfig(Config{Users: []User{user}}); err != nil {
			return err
		}
		return nil
	}

	config, err := ReadConfig()
	if err != nil {
		return err
	}

	if config.Users == nil {
		config.Users = []User{user}
	} else {
		err := CheckUser(user, config)
		if err != nil {
			return err
		}

		err = CheckAlias(user.Alias, config.Users)
		if err != nil {
			return err
		}

		config.Users = append(config.Users, user)
	}

	return CreateConfig(config)
}

func RemoveUser(idx int, users []User) error {
	newUsers := make([]User, len(users)-1)
	if idx+1 == len(users) {
		newUsers = users[:idx]
	} else {
		newUsers = append(users[:idx], users[idx+1:]...)
	}

	return CreateConfig(Config{Users: newUsers})
}

func ModifyUser(idx int, newUser User) error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	if newUser.Name != "" {
		config.Users[idx].Name = newUser.Name
	}

	if newUser.Email != "" {
		config.Users[idx].Email = newUser.Email
	}

	if newUser.GpgKeyID != "" {
		config.Users[idx].GpgKeyID = newUser.GpgKeyID
	}

	err = CheckAlias(newUser.Alias, config.Users)
	if err != nil {
		return err
	}

	config.Users[idx].Alias = newUser.Alias
	return CreateConfig(config)
}

func CheckUser(newUser User, config Config) error {
	for _, user := range config.Users {
		if user.Email == newUser.Email && user.Name == newUser.Name {
			msg := fmt.Sprintf("User %s <%s> already exists", user.Name, user.Email)
			return errors.New(msg)
		}
	}

	return nil
}

func CheckAlias(alias string, users []User) error {
	if alias == "" {
		return nil
	}

	for _, user := range users {
		if user.Alias == alias {
			msg := fmt.Sprintf("A user with alias [%s] already exists: %s <%s>", alias, user.Name, user.Email)
			return errors.New(msg)
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	if !govalidator.IsExistingEmail(email) {
		return errors.New("Invalid email address")
	}

	return nil
}

func ValidateModifiedEmail(email string) error {
	if email == "" {
		return nil
	}

	if !govalidator.IsExistingEmail(email) {
		return errors.New("Invalid email address")
	}

	return nil
}
