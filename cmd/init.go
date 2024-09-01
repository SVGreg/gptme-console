/*
Copyright © 2024 GPTMe
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SVGreg/gptme-console/gpt"
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
	initCmd.PersistentFlags().StringP("path", "p", "", "Path where to store configuration file")
}

func initRun(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	if key != "" {
		fmt.Println("Specified key is", key)
	} else {
		cmd.Help()
		os.Exit(0)
	}

	if len(args) > 0 {
		cmd.Help()
		os.Exit(0)
	}

	path, _ := cmd.Flags().GetString("path")
	if path == "" {
		path = ".gptme-config.json"
	} else {
		path += string(os.PathSeparator) + ".gptme-config.json"
	}
	fmt.Println("Config path is", path)

	err := SaveConfig(path, gpt.Config{APIKey: key})
	if err != nil {
		_ = fmt.Errorf("%v", err)
		os.Exit(0)
	}
}

func SaveConfig(filename string, config gpt.Config) error {
	res, err := json.Marshal(config)
	if err != nil {
		return err
	}

	werr := os.WriteFile(filename, res, 0644)
	if werr != nil {
		return werr
	}

	fmt.Println("Configuration", string(res), "saved at", filename)

	return nil
}
