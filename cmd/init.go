/*
Copyright © 2024 GPTMe
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/SVGreg/gptme-console/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises GPT environment",
	Long:  `Initialises GPT environment and configures API key`,
	Run:   initRun,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringP("key", "k", "", "Specifies API key to work with")
	initCmd.PersistentFlags().StringP("org", "o", "", "OpenAI Organization Id")
	initCmd.PersistentFlags().StringP("project", "j", "", "OpenAI Project Id")
}

func initRun(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Help()
		os.Exit(0)
	}

	path, _ := cmd.Flags().GetString("path")
	path = config.MakePath(path)
	fmt.Println("Config path is", path)

	key, _ := cmd.Flags().GetString("key")
	if key == "" {
		cmd.Help()
		os.Exit(1)
	}

	org, _ := cmd.Flags().GetString("org")
	if org == "" {
		cmd.Help()
		os.Exit(1)
	}

	proj, _ := cmd.Flags().GetString("project")
	if proj == "" {
		cmd.Help()
		os.Exit(1)
	}

	err := config.Save(path, config.Config{
		OrganizationId: org, 
		ProjectId: proj, 
		APIKey: key,
	})
	if err != nil {
		_ = fmt.Errorf("%v", err)
		os.Exit(0)
	}
}
