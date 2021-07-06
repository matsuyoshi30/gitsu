package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matsuyoshi30/gitsu/internal/constants"
	"github.com/matsuyoshi30/gitsu/internal/models"
	"github.com/matsuyoshi30/gitsu/internal/utils"
)

var (
	ErrConfigFileDoesNotExist = errors.New("Config file does not exist")
	ErrUserIndexOutOfBounds   = errors.New("User index out of bounds")
	ErrNoDefaultUser          = errors.New("No default user")
	ErrNoUserWithAlias        = errors.New("No user with this alias")
)

// Config describes the structure of the JSON based config file
type Config struct {
	Version string        `json:"version"`
	Users   []models.User `json:"users"`
}

// Dir returns the config directory
func Dir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, constants.ConfigDir), nil
}

// Path returns the config file path
func Path() (string, error) {
	configDir, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, constants.ConfigFileName), nil
}

// Exists returns if the config file exists
func Exists() (bool, error) {
	configFilePath, err := Path()
	if err != nil {
		return false, err
	}

	return utils.FileExists(configFilePath), nil
}

// Read reads the config file. Returns the config as a struct or an error if failed to read
func Read() (*Config, error) {
	exists, err := Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrConfigFileDoesNotExist
	}

	configFilePath, err := Path()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Write writes config data to the config file. Returns an error if failed to write data
func Write(c *Config) error {
	configDir, err := Dir()
	if err != nil {
		return err
	}

	// Check if the 'git-su' directory exists. If not, create it
	dirExists := utils.DirExists(configDir)
	if !dirExists {
		err := os.Mkdir(configDir, 0744)
		if err != nil {
			return err
		}
	}

	configFilePath, err := Path()
	if err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(c, "", constants.JsonIndent)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// AddUser adds a new user to the config or returns an error if new user is invalid
func (c *Config) AddUser(user *models.User) error {
	err := c.isValidUser(user)
	if err != nil {
		return err
	}

	c.Users = append(c.Users, *user)
	return nil
}

// ModifyUser modifies an existing user in the config or returns an error if the index is out of bounds
func (c *Config) ModifyUser(index int, modifiedUser *models.User) error {
	if index < 0 || index > len(c.Users)-1 {
		return ErrUserIndexOutOfBounds
	}

	user := c.Users[index]
	user.Modify(
		modifiedUser.Name,
		modifiedUser.Email,
		modifiedUser.Alias,
		modifiedUser.GpgKeyID,
	)

	err := c.isValidUser(&user)
	if err != nil {
		return err
	}

	c.Users[index] = user
	return nil
}

// DeleteUser deletes a user from the config or returns an error if the index is out of bounds
func (c *Config) DeleteUser(index int) error {
	if index < 0 || index > len(c.Users)-1 {
		return ErrUserIndexOutOfBounds
	}

	c.Users[index] = c.Users[len(c.Users)-1]
	c.Users = c.Users[:len(c.Users)-1]

	return nil
}

// SelectUser selects an existing user from the config or returns an error if the index is out of bounds
func (c *Config) SelectUser(index int) (*models.User, error) {
	if index < 0 || index > len(c.Users)-1 {
		return nil, ErrUserIndexOutOfBounds
	}

	return &c.Users[index], nil
}

// SelectDefaultUser returns the default user or nil if there is no default user
func (c *Config) SelectDefaultUser() (*models.User, error) {
	for _, user := range c.Users {
		if user.Alias == constants.DefaultAlias {
			return &user, nil
		}
	}
	return nil, ErrNoDefaultUser
}

// SelectUserByAlias returns a user by alias or nil if there is no such user
func (c *Config) SelectUserByAlias(alias string) (*models.User, error) {
	for _, user := range c.Users {
		if user.Alias == alias {
			return &user, nil
		}
	}
	return nil, ErrNoUserWithAlias
}

// UserList returns a list (slice) of formatted user data
func (c *Config) UserList() []string {
	var padding int = 0
	var list []string
	for _, user := range c.Users {
		if len(user.Alias) > padding {
			padding = len(user.Alias) + 2
		}
	}
	for _, user := range c.Users {
		list = append(list, user.Format(padding))
	}
	return list
}

// Reset resets the saved user profiles
func (c *Config) Reset() {
	c.Users = []models.User{}
}

// isValidUser returns if the provided user is valid
func (c *Config) isValidUser(newUser *models.User) error {
	for _, user := range c.Users {
		if user.Name == newUser.Name && user.Email == newUser.Email {
			return fmt.Errorf("User %s <%s> already exists", user.Name, user.Email)
		}

		if newUser.Alias != "" && user.Alias == newUser.Alias {
			return fmt.Errorf(
				"A user with alias [%s] already exists: %s <%s>",
				newUser.Alias,
				newUser.Name,
				newUser.Email,
			)
		}
	}
	return nil
}
