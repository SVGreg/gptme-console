/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package cmd

import (
	"fmt"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/internal/logger"
	"github.com/SVGreg/gptme-console/internal/validation"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long: `Initialize the configuration file with API key and other settings.
This command creates a configuration file in the specified directory or in the home directory by default.`,
	Run: initRun,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringP("key", "k", "", "OpenAI API key (required)")
	initCmd.PersistentFlags().StringP("org", "o", "", "OpenAI Organization ID (required)")
	initCmd.PersistentFlags().StringP("project", "j", "", "OpenAI Project ID (required)")

	// Mark required flags
	initCmd.MarkPersistentFlagRequired("key")
	initCmd.MarkPersistentFlagRequired("org")
	initCmd.MarkPersistentFlagRequired("project")
}

func initRun(cmd *cobra.Command, args []string) {
	// Get flag values
	apiKey, _ := cmd.Flags().GetString("key")
	orgID, _ := cmd.Flags().GetString("org")
	projectID, _ := cmd.Flags().GetString("project")
	path, _ := cmd.Flags().GetString("path")

	// Validate inputs
	if err := validation.ValidateAPIKey(apiKey); err != nil {
		logger.Fatal("Invalid API key: %v", err)
	}

	if err := validation.ValidateOrganizationID(orgID); err != nil {
		logger.Fatal("Invalid organization ID: %v", err)
	}

	if err := validation.ValidateProjectID(projectID); err != nil {
		logger.Fatal("Invalid project ID: %v", err)
	}

	// Create config
	cfg := config.Config{
		APIKey:         apiKey,
		OrganizationId: orgID,
		ProjectId:      projectID,
	}

	// Save config
	configPath := config.MakePath(path)
	if err := config.Save(configPath, cfg); err != nil {
		logger.Fatal("Failed to save configuration: %v", err)
	}

	fmt.Printf("Configuration saved successfully to %s\n", configPath)
}
