// FILE: cmd/ask_test.go
package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/gpt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAskRun(t *testing.T) {
	// Mock the dependencies
	config.Read = func(path string) (config.Config, error) {
		return config.Config{APIKey: "test-api-key"}, nil
	}
	gpt.Request = func(question string, config config.Config) string {
		return "This is a test response."
	}

	// Capture the output
	var output bytes.Buffer
	log.SetOutput(&output)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Too many arguments",
			args:     strings.Split(strings.Repeat("word ", 31), " "),
			expected: "Usage:\n",
		},
		{
			name:     "Valid question",
			args:     []string{"What", "is", "the", "capital", "of", "France?"},
			expected: "Q: What is the capital of France?\nA: This is a test response.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("path", "test-path", "config file path")

			// Reset the output buffer
			output.Reset()

			// Run the askRun function
			askRun(cmd, tt.args)

			// Check the output
			assert.Contains(t, output.String(), tt.expected)
		})
	}
}
