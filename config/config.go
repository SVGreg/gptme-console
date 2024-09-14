/*
Copyright © 2024 GPTMe
*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const Filename string = ".gptme-config.json"

type Config struct {
	OrganizationId string
	ProjectId string
	APIKey string
}

func MakePath(path string) string {
	if path == "" {
		return Filename
	} else {
		return path + string(os.PathSeparator) + Filename
	}
}

func Save(filename string, config Config) error {
	res, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, res, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Configuration", string(res), "saved at", filename)

	return nil
}

func Read(filename string) (Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
