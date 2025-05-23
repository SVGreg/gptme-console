/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/spf13/cobra"

	markdown "github.com/MichaelMure/go-term-markdown"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Asks GPT you question",
	Long: `Asks GPT you question. 
This command sends a question to GPT and displays the response in markdown format.
This is ad-hoc question&answer and will not be stored in session.`,
	Run: askRun,
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func askRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Please provide a question")
		cmd.Help()
		return
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
	fmt.Println("A:", string(markdown.Render(response, 120, 2)))
}
