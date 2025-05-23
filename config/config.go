/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/errors"
	"github.com/SVGreg/gptme-console/internal/logger"
	"github.com/SVGreg/gptme-console/internal/validation"
)

// ReadFunc is a type for the Read function
type ReadFunc func(filename string) (Config, error)

// DefaultRead is the default implementation of Read
var DefaultRead ReadFunc = func(filename string) (Config, error) {
	logger.Debug("Attempting to read config from: %s", filename)

	// Try to read from file first
	config, err := readFromFile(filename)
	if err == nil {
		logger.ConfigLoaded(filename)
		return config, nil
	}

	// If file doesn't exist, try environment variables
	logger.Debug("Config file not found, trying environment variables")
	config, envErr := readFromEnv()
	if envErr == nil {
		logger.Info("Config loaded from environment variables")
		return config, nil
	}

	// If both fail, return the file error
	logger.ConfigNotFound(filename)
	return Config{}, errors.NewConfigError("read", err)
}

// Read is the function that can be reassigned for testing
var Read = DefaultRead

type Config struct {
	OrganizationId string `json:"OrganizationId"`
	ProjectId      string `json:"ProjectId"`
	APIKey         string `json:"APIKey"`
}

func readFromFile(filename string) (Config, error) {
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

func readFromEnv() (Config, error) {
	config := Config{
		APIKey:         os.Getenv(constants.EnvAPIKey),
		OrganizationId: os.Getenv(constants.EnvOrgID),
		ProjectId:      os.Getenv(constants.EnvProjectID),
	}

	// Validate that at least API key is present
	if config.APIKey == "" {
		return Config{}, errors.NewConfigError("env",
			errors.NewValidationError("API key", "",
				fmt.Errorf("environment variable %s not set", constants.EnvAPIKey)))
	}

	return config, nil
}

func MakePath(path string) string {
	// Check environment variable first
	if envPath := os.Getenv(constants.EnvConfigPath); envPath != "" {
		return filepath.Join(envPath, constants.ConfigFileName)
	}

	if path == "" {
		// Get user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory if home directory cannot be determined
			home = "."
		}
		return filepath.Join(home, constants.ConfigFileName)
	}
	return filepath.Join(path, constants.ConfigFileName)
}

func Save(filename string, config Config) error {
	logger.Debug("Saving config to: %s", filename)

	// Validate config before saving
	if err := validateConfig(config); err != nil {
		return errors.NewConfigError("validation", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, constants.DirPerms); err != nil {
		return errors.NewConfigError("mkdir", err)
	}

	res, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.NewConfigError("marshal", err)
	}

	err = os.WriteFile(filename, res, constants.FilePerms)
	if err != nil {
		return errors.NewConfigError("write", err)
	}

	logger.ConfigSaved(filename)
	return nil
}

func validateConfig(config Config) error {
	if err := validation.ValidateAPIKey(config.APIKey); err != nil {
		return err
	}

	if err := validation.ValidateOrganizationID(config.OrganizationId); err != nil {
		return err
	}

	if err := validation.ValidateProjectID(config.ProjectId); err != nil {
		return err
	}

	return nil
}
