package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Users []User `json:"Users"`
}

const ConfigFile = "config.json"

func ConfigDir() string {
	configDirName := "gitsu-go"

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "/Library/Preferences", configDirName)
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "gitsu-go")
	default:
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), configDirName)
		}
		return filepath.Join(os.Getenv("HOME"), "/.config", configDirName)
	}
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), ConfigFile)
}

func IsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func CreateConfig(config Config) error {
	if !IsExist(ConfigDir()) {
		if err := os.Mkdir(ConfigDir(), 0744); err != nil {
			return err
		}
	}

	file, err := os.Create(ConfigPath())
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

	data, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
