package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
)

var CONFIG_PATHS = []string{
	".config/go-git/config.json",
}

type Config struct {
	SCM_Name string `json:"scm_name"`
	Token    string `json:"token"`
}

func ReadConfig() (config Config, err error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	var configFilePath string
	for _, path := range CONFIG_PATHS {
		fullPath := fmt.Sprintf("%s/%s", userHomeDir, path)
		_, err := os.Stat(fullPath)
		if err == nil{
			slog.Debug("File Path Exists", "fullPath", fullPath)
			configFilePath = fullPath
		} else if os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", fullPath)
		} else {
			fmt.Printf("Error checking file: %s\n", fullPath)
		}
	}

	slog.Debug("Config File Path", "configFilePath", configFilePath)

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	if config.SCM_Name == "" {
		log.Println("SCM Type is not set in the config file.")
	} else {
		slog.Debug("SCM Type", "scm_name", config.SCM_Name)
	}

	return config, err
}