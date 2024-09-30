/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Asks GPT you question",
	Long: `Asks GPT you question. 
Please type the question right after 'ask' command. It is limited to 20 words.`,
	Run: askRun,
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func askRun(cmd *cobra.Command, args []string) {
	if len(args) > 20 {
		cmd.Help()
		os.Exit(0)
	}

	question := strings.Join(args, " ")
	fmt.Println("Q:", question)

	// Read config: api key
	path, _ := cmd.Flags().GetString("path")
	config, err := config.Read(config.MakePath(path))
	if err != nil {
		log.Fatalln("Unable to read config", err)
	}

	// Request answer
	response := gpt.Request(question, config)
	fmt.Println("A:", response)
}
