/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Asks GPT you question",
	Long: `Asks GPT you question. 
Please specify exactly one string parameter with the requested question.`,
	Run: askRun,
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func askRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(0)
	}

	fmt.Println("Q:", args[0])
	for arg := range args {
		fmt.Println(arg)
	}

	// Read config: api key
	config, err := config.Read(config.MakePath(""))
	if err != nil {
		log.Fatalln("Unable to read config", err)
	}

	// Request answer
	gpt.Request(args[0], config)

}

