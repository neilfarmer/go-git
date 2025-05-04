package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/neilfarmer/go-git/internal/config"
	"github.com/neilfarmer/go-git/internal/github"
)

var CONFIG_PATHS = []string{
	".config/go-git/config.json",
}

func main() {
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
			fmt.Printf("File exists: %s\n", fullPath)
			configFilePath = fullPath
		} else if os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", fullPath)
		} else {
			fmt.Printf("Error checking file: %s\n", fullPath)
		}
	}

	fmt.Printf("The full path to the config file is: %s\n", configFilePath)

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config config.Config
	json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Printf("SCM Type: %s\n", config.SCM_Name)
	if config.SCM_Name == "" {
		fmt.Println("SCM Type is not set in the config file.")
	} else {
		fmt.Printf("SCM Type is set to: %s\n", config.SCM_Name)
	}

	if config.SCM_Name == "github" {
		fmt.Println("GitHub configuration detected.")
		github.GetRepos(config)
	}
}
