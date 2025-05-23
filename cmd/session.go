/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/logger"
	"github.com/SVGreg/gptme-console/internal/validation"
	"github.com/SVGreg/gptme-console/session"
	"github.com/spf13/cobra"

	markdown "github.com/MichaelMure/go-term-markdown"
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
	sessionCmd.PersistentFlags().StringP("description", "d", "", "Description for the new session (used with --start)")
	sessionCmd.PersistentFlags().StringArrayP("tags", "t", []string{}, "Tags for the new session (used with --start)")
	sessionCmd.PersistentFlags().StringP("use", "u", "", "Uses session with specified name. Raises error if session does not exist.")
	sessionCmd.PersistentFlags().BoolP("list", "l", false, "Prints the list of stored sessions")
	sessionCmd.PersistentFlags().StringP("cat", "c", "", "Prints history of specified session. If empty, shows current session.")
	sessionCmd.PersistentFlags().Bool("clean", false, "Cleans up all stored sessions")
	sessionCmd.PersistentFlags().StringP("ask", "a", "", "Asks a question in the current session")
}

func sessionRun(cmd *cobra.Command, args []string) {
	sm, err := session.NewSessionManager()
	if err != nil {
		logger.Fatal("Error initializing session manager: %v", err)
	}

	// Handle list flag
	if list, _ := cmd.Flags().GetBool("list"); list {
		handleListSessions(sm)
		return
	}

	// Handle cat flag
	if cmd.Flags().Changed("cat") {
		catName, _ := cmd.Flags().GetString("cat")
		handleCatSession(sm, catName)
		return
	}

	// Handle clean flag
	if clean, _ := cmd.Flags().GetBool("clean"); clean {
		handleCleanSessions(sm)
		return
	}

	// Handle start flag
	if startName, _ := cmd.Flags().GetString("start"); startName != "" {
		description, _ := cmd.Flags().GetString("description")
		tags, _ := cmd.Flags().GetStringArray("tags")
		handleStartSession(sm, startName, description, tags)
		return
	}

	// Handle use flag
	if useName, _ := cmd.Flags().GetString("use"); useName != "" {
		handleUseSession(sm, useName)
		return
	}

	// Handle ask flag
	if question, _ := cmd.Flags().GetString("ask"); question != "" {
		handleAskInSession(sm, question, cmd)
		return
	}

	// If no flags are set, show help
	cmd.Help()
}

func handleListSessions(sm *session.SessionManager) {
	sessions, err := sm.ListSessions()
	if err != nil {
		logger.Fatal("Error listing sessions: %v", err)
	}
	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		return
	}
	fmt.Println("Available sessions:")
	for _, s := range sessions {
		currentMarker := ""
		if s.Name == sm.Current {
			currentMarker = " (current)"
		}

		tagsStr := ""
		if len(s.Tags) > 0 {
			tagsStr = fmt.Sprintf(" [%s]", strings.Join(s.Tags, ", "))
		}

		descStr := ""
		if s.Description != "" {
			descStr = fmt.Sprintf(" - %s", s.Description)
		}

		fmt.Printf("Session '%s': Created at %s, updated at %s, %d messages%s%s%s\n",
			s.Name,
			s.CreatedAt.Format("2006-01-02 15:04:05"),
			s.UpdatedAt.Format("2006-01-02 15:04:05"),
			len(s.Messages),
			descStr,
			tagsStr,
			currentMarker)
	}
}

func handleCatSession(sm *session.SessionManager, sessionName string) {
	if sessionName == "" {
		if sm.Current == "" {
			logger.Fatal("No current session. Use --start or --use to select a session first.")
		}
		sessionName = sm.Current
	}

	s := sm.GetSession(sessionName)
	if s == nil {
		logger.Fatal("Session '%s' not found", sessionName)
	}

	fmt.Printf("Session '%s' history:\n", sessionName)
	if s.Description != "" {
		fmt.Printf("Description: %s\n", s.Description)
	}
	if len(s.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(s.Tags, ", "))
	}
	fmt.Println()

	for _, msg := range s.Messages {
		fmt.Printf("[%s] %s: %s\n",
			msg.Timestamp.Format("2006-01-02 15:04:05"),
			msg.Role,
			string(markdown.Render(msg.Content, constants.MarkdownWidth, constants.MarkdownIndent)))
	}
}

func handleCleanSessions(sm *session.SessionManager) {
	if err := sm.Clean(); err != nil {
		logger.Fatal("Error cleaning sessions: %v", err)
	}
	fmt.Println("All sessions cleaned")
}

func handleStartSession(sm *session.SessionManager, name, description string, tags []string) {
	s, err := sm.CreateSession(name, description, tags)
	if err != nil {
		logger.Fatal("Error: %v", err)
	}
	fmt.Printf("Created and using session '%s'\n", s.Name)
}

func handleUseSession(sm *session.SessionManager, name string) {
	if err := sm.SetCurrent(name); err != nil {
		logger.Fatal("Error: %v", err)
	}
	fmt.Printf("Using session '%s'\n", name)
}

func handleAskInSession(sm *session.SessionManager, question string, cmd *cobra.Command) {
	if sm.Current == "" {
		logger.Fatal("No current session. Use --start or --use to select a session first.")
	}

	// Validate question
	if err := validation.ValidateQuestion(question); err != nil {
		logger.Fatal("Error: %v", err)
	}

	fmt.Println("Q:", question)

	// Read config
	path, _ := cmd.Flags().GetString("path")
	config, err := config.Read(config.MakePath(path))
	if err != nil {
		logger.Fatal("Unable to read config: %v", err)
	}

	// Store the question in the session
	if err := sm.AddMessage(sm.Current, "user", question); err != nil {
		logger.Fatal("Error storing question: %v", err)
	}

	// Get conversation history for context
	history, err := sm.GetConversationHistory(sm.Current)
	if err != nil {
		logger.Fatal("Error getting conversation history: %v", err)
	}

	// Convert session messages to GPT context (excluding the last message which is the current question)
	var context []gpt.ContextMessage
	for i, msg := range history {
		if i < len(history)-1 { // Skip the last message (current question)
			context = append(context, gpt.ContextMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// Request answer with context
	response := gpt.RequestWithContext(question, config, context)
	fmt.Println("A:", string(markdown.Render(response, constants.MarkdownWidth, constants.MarkdownIndent)))

	// Store the response in the session
	if err := sm.AddMessage(sm.Current, "assistant", response); err != nil {
		logger.Fatal("Error storing response: %v", err)
	}
}
