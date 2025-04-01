/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/SVGreg/gptme-console/session"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Runs in session mode and manages modes.",
	Long:  `Sessions allows to store communication history and use it as context for further questions.`,
	Run:   sessionRun,
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.PersistentFlags().StringP("start", "s", "", "Creates new session with specified name and makes it current")
	sessionCmd.PersistentFlags().StringP("use", "u", "", "Uses session with specified name. Raises error if session does not exist.")
	sessionCmd.PersistentFlags().BoolP("list", "l", false, "Prints the list of stored sessions")
	sessionCmd.PersistentFlags().StringP("cat", "c", "", "Prints history of specified session. Recommended to use with 'less' or 'more'.")
	sessionCmd.PersistentFlags().Bool("clean", false, "Cleans up all stored sessions")
}

func sessionRun(cmd *cobra.Command, args []string) {
	sm, err := session.NewSessionManager()
	if err != nil {
		fmt.Printf("Error initializing session manager: %v\n", err)
		os.Exit(1)
	}

	// Handle list flag
	if list, _ := cmd.Flags().GetBool("list"); list {
		sessions, err := sm.ListSessions()
		if err != nil {
			fmt.Printf("Error listing sessions: %v\n", err)
			os.Exit(1)
		}
		if len(sessions) == 0 {
			fmt.Println("No sessions found")
			return
		}
		fmt.Println("Available sessions:")
		for _, s := range sessions {
			fmt.Printf("Session '%s': Created at %s, %d messages\n", s.Name, s.CreatedAt.Format("2006-01-02 15:04:05"), len(s.Messages))
		}
		return
	}

	// Handle cat flag
	if catName, _ := cmd.Flags().GetString("cat"); catName != "" {
		s := sm.GetSession(catName)
		if s == nil {
			fmt.Printf("Session '%s' not found\n", catName)
			return
		}
		fmt.Printf("Session '%s' history:\n", catName)
		for _, msg := range s.Messages {
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Role, msg.Content)
		}
		return
	}

	// Handle clean flag
	if clean, _ := cmd.Flags().GetBool("clean"); clean {
		if err := sm.Clean(); err != nil {
			fmt.Printf("Error cleaning sessions: %v\n", err)
			return
		}
		fmt.Println("All sessions cleaned")
		return
	}

	// Handle start flag
	if startName, _ := cmd.Flags().GetString("start"); startName != "" {
		s, err := sm.CreateSession(startName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created and using session '%s'\n", s.Name)
		return
	}

	// Handle use flag
	if useName, _ := cmd.Flags().GetString("use"); useName != "" {
		if err := sm.SetCurrent(useName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Using session '%s'\n", useName)
		return
	}

	// If no flags are set, show help
	cmd.Help()
}
