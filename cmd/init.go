/*
Copyright © 2024 GPTMe
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/SVGreg/gptme-console/config"
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
}

func initRun(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	if key == "" {
		cmd.Help()
		os.Exit(0)
	}

	if len(args) > 0 {
		cmd.Help()
		os.Exit(0)
	}

	path, _ := cmd.Flags().GetString("path")
	path = config.MakePath(path)
	fmt.Println("Config path is", path)

	err := config.Save(path, config.Config{APIKey: key})
	if err != nil {
		_ = fmt.Errorf("%v", err)
		os.Exit(0)
	}
}
