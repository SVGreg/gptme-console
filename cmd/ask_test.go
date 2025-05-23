// FILE: cmd/ask_test.go
package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/SVGreg/gptme-console/internal/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAskRun(t *testing.T) {
	// Mock the dependencies
	originalRead := config.Read
	originalRequest := gpt.Request

	defer func() {
		config.Read = originalRead
		gpt.Request = originalRequest
	}()

	config.Read = func(path string) (config.Config, error) {
		return config.Config{APIKey: "test-api-key"}, nil
	}
	gpt.Request = func(question string, config config.Config) string {
		return "This is a test response."
	}

	// Capture the output
	var output bytes.Buffer
	logger.SetLevel(logger.LevelError) // Suppress lower level logs for testing

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "No arguments",
			args:     []string{},
			expected: "Error: Please provide a question",
		},
		{
			name:     "Valid question",
			args:     []string{"What", "is", "the", "capital", "of", "France?"},
			expected: "Q: What is the capital of France?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("path", "test-path", "config file path")

			// Reset the output buffer
			output.Reset()

			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the askRun function
			askRun(cmd, tt.args)

			// Restore stdout and get output
			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output.Write(buf[:n])

			// Check the output contains expected text
			assert.Contains(t, output.String(), tt.expected)
		})
	}
}
