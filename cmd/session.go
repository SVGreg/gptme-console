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

	sessionCmd.PersistentFlags().IntP("run", "r", -1, "Runs session with specified Id or starts new session if no Id set")
	sessionCmd.PersistentFlags().BoolP("list", "l", false, "Prints the list of stored sessions")
	sessionCmd.PersistentFlags().IntP("cat", "c", -1, "Prints history of specified session. Recommended to use with 'less' or 'more'.")
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
		sessions := sm.ListSessions()
		if len(sessions) == 0 {
			fmt.Println("No sessions found")
			return
		}
		fmt.Println("Available sessions:")
		for _, s := range sessions {
			fmt.Printf("Session %d: Created at %s, %d messages\n", s.ID, s.CreatedAt.Format("2006-01-02 15:04:05"), len(s.Messages))
		}
		return
	}

	// Handle cat flag
	if catID, _ := cmd.Flags().GetInt("cat"); catID != -1 {
		s := sm.GetSession(catID)
		if s == nil {
			fmt.Printf("Session %d not found\n", catID)
			return
		}
		fmt.Printf("Session %d history:\n", catID)
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

	// Handle run flag
	if runID, _ := cmd.Flags().GetInt("run"); runID != -1 {
		if runID == 0 {
			s := sm.CreateSession()
			fmt.Printf("Created new session %d\n", s.ID)
		} else {
			s := sm.GetSession(runID)
			if s == nil {
				fmt.Printf("Session %d not found\n", runID)
				return
			}
			fmt.Printf("Using session %d\n", runID)
		}
		// TODO: Implement interactive session mode
		fmt.Println("Interactive session mode not implemented yet")
		return
	}

	// If no flags are set, show help
	cmd.Help()
}
