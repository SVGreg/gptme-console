/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/errors"
)

var (
	// Session name should be alphanumeric with hyphens and underscores
	sessionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	// API key patterns (OpenAI format)
	apiKeyRegex = regexp.MustCompile(`^sk-[a-zA-Z0-9-_]+$`)

	// Organization ID pattern
	orgIDRegex = regexp.MustCompile(`^org-[a-zA-Z0-9]+$`)

	// Project ID pattern
	projectIDRegex = regexp.MustCompile(`^proj_[a-zA-Z0-9]+$`)
)

// ValidateSessionName validates a session name
func ValidateSessionName(name string) error {
	if name == "" {
		return errors.NewValidationError("session name", name, fmt.Errorf("cannot be empty"))
	}

	if len(name) > 50 {
		return errors.NewValidationError("session name", name, fmt.Errorf("cannot be longer than 50 characters"))
	}

	if !sessionNameRegex.MatchString(name) {
		return errors.NewValidationError("session name", name, fmt.Errorf("can only contain letters, numbers, hyphens, and underscores"))
	}

	// Reserved names
	reserved := []string{"sessions", "config", "current", "list", "all"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return errors.NewValidationError("session name", name, fmt.Errorf("'%s' is a reserved name", r))
		}
	}

	return nil
}

// ValidateAPIKey validates an OpenAI API key
func ValidateAPIKey(key string) error {
	if key == "" {
		return errors.NewValidationError("API key", key, fmt.Errorf("cannot be empty"))
	}

	if !apiKeyRegex.MatchString(key) {
		return errors.NewValidationError("API key", key, fmt.Errorf("invalid format (should start with 'sk-')"))
	}

	return nil
}

// ValidateOrganizationID validates an OpenAI organization ID
func ValidateOrganizationID(orgID string) error {
	if orgID == "" {
		return errors.NewValidationError("organization ID", orgID, fmt.Errorf("cannot be empty"))
	}

	if !orgIDRegex.MatchString(orgID) {
		return errors.NewValidationError("organization ID", orgID, fmt.Errorf("invalid format (should start with 'org-')"))
	}

	return nil
}

// ValidateProjectID validates an OpenAI project ID
func ValidateProjectID(projectID string) error {
	if projectID == "" {
		return errors.NewValidationError("project ID", projectID, fmt.Errorf("cannot be empty"))
	}

	if !projectIDRegex.MatchString(projectID) {
		return errors.NewValidationError("project ID", projectID, fmt.Errorf("invalid format (should start with 'proj_')"))
	}

	return nil
}

// ValidateQuestion validates a question for word limit
func ValidateQuestion(question string) error {
	if question == "" {
		return errors.NewValidationError("question", question, fmt.Errorf("cannot be empty"))
	}

	words := strings.Fields(question)
	if len(words) > constants.MaxQuestionWords {
		return errors.NewValidationError("question", question, fmt.Errorf("cannot exceed %d words (got %d)", constants.MaxQuestionWords, len(words)))
	}

	return nil
}
