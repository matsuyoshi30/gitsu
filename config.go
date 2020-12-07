package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Users []User `json:"Users"`
}

const ConfigFile = "config.json"

func ConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "gitsu-go"), nil
}

func ConfigPath() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, ConfigFile), nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func CreateConfig(config Config) error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}
	if !IsExist(configDir) {
		if err := os.Mkdir(configDir, 0744); err != nil {
			return err
		}
	}

	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func ReadConfig() (Config, error) {
	var config Config

	configPath, err := ConfigPath()
	if err != nil {
		return config, err
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
