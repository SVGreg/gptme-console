/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/logger"
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

	// Read config
	path, _ := cmd.Flags().GetString("path")
	config, err := config.Read(config.MakePath(path))
	if err != nil {
		logger.Fatal("Unable to read config: %v", err)
	}

	// Request answer (no context for standalone ask)
	response := gpt.Request(question, config)
	fmt.Println("A:", string(markdown.Render(response, constants.MarkdownWidth, constants.MarkdownIndent)))
}
