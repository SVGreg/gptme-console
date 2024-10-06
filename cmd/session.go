/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Runs in session mode and manages modes.",
	Long:  `Sessions allows to store communication history and use it as context for further questions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("session called")
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.PersistentFlags().IntP("run", "r", -1, "Runs session with specified Id or starts new session if no Id set")
	sessionCmd.PersistentFlags().BoolP("list", "l", false, "Prints the list of stored sessions")
	sessionCmd.PersistentFlags().IntP("print", "p", -1, "Prints history of specified session. Recommended to use with 'less' or 'more'.")
	sessionCmd.PersistentFlags().Bool("clean", false, "Cleans up all stored sessions")
}
