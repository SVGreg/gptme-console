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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gptme-console",
	Short: "Minimal console client to ChatGPT API",
	Long: `Minimal console client to ChatGPT API
Require to initialization ahead of use. 
Please specify API Key using init command. 
Specify optional path to config file.
	`,
	Run: rootRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", "", "Path where to store configuration file")
}

func rootRun(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("path")
	path = config.MakePath(path)

	config, err := config.Read(path)
	if err != nil {
		fmt.Println("Can't find config at path: ", path)
		cmd.Help()
		os.Exit(0)
	}

	fmt.Println("Found config at path : ", path)
	fmt.Printf("\t- Organization Id : %s...\n", config.OrganizationId[0:16])
	fmt.Printf("\t- Project Id      : %s...\n", config.ProjectId[0:16])
	fmt.Printf("\t- API key         : %s...\n", config.APIKey[0:16])
}
