/*
Copyright © 2024 GPTMe
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/SVGreg/gptme-console/gpt"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises GPT environment",
	Long:  `Initialises GPT environment and configures API key`,
	Run:   addRun,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addRun(cmd *cobra.Command, args []string) {
	fmt.Println("Please specify API key")
	for _, arg := range args {
		fmt.Println(arg)
	}
}
