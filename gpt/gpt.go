/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/

package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SVGreg/gptme-console/config"
	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/errors"
	"github.com/SVGreg/gptme-console/internal/logger"
)

// RequestFunc is a type for the Request function
type RequestFunc func(question string, config config.Config) string

// RequestWithContextFunc is a type for the RequestWithContext function
type RequestWithContextFunc func(question string, config config.Config, context []ContextMessage) string

// ContextMessage represents a message in conversation context
type ContextMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DefaultRequest is the default implementation of Request (without context)
var DefaultRequest RequestFunc = func(question string, config config.Config) string {
	return DefaultRequestWithContext(question, config, nil)
}

// DefaultRequestWithContext is the default implementation with conversation context
var DefaultRequestWithContext RequestWithContextFunc = func(question string, config config.Config, context []ContextMessage) string {
	logger.GPTRequest(question)

	// Build messages array with context
	messages := make([]map[string]string, 0, len(context)+1)

	// Add conversation history
	for _, msg := range context {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// Add current question
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": question,
	})

	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		logger.GPTError(errors.NewGPTError("marshal_messages", err))
		return "[ERROR] Failed to prepare request"
	}

	bodyString := fmt.Sprintf(`{
		"model": "%s",
		"messages": %s,
		"temperature": %g
	}`, constants.GPTModel, string(messagesJSON), constants.GPTTemperature)

	reader := strings.NewReader(bodyString)

	request, err := http.NewRequest(http.MethodPost, constants.CompletionURL, reader)
	if err != nil {
		logger.GPTError(errors.NewGPTError("create_request", err))
		return "[ERROR] Can't create request"
	}

	request.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)
	request.Header.Set(constants.HeaderOrganization, config.OrganizationId)
	request.Header.Set(constants.HeaderProject, config.ProjectId)
	request.Header.Set(constants.HeaderAuth, constants.HeaderAuthPrefix+config.APIKey)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.GPTError(errors.NewGPTError("http_request", err))
		return fmt.Sprintf("[ERROR] HTTP error: %v", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.GPTError(errors.NewGPTError("read_response", err))
		return "[ERROR] Failed to read response"
	}

	var gptError GPTError
	if err := json.Unmarshal(resBody, &gptError); err == nil && len(gptError.Error.Message) > 0 {
		logger.GPTError(errors.NewGPTError("api_error", fmt.Errorf(gptError.Error.Message)))
		return fmt.Sprintf("[ERROR] GPT reports error: %s", gptError.Error.Message)
	}

	var gptContent GPTResponse
	if err := json.Unmarshal(resBody, &gptContent); err != nil {
		logger.GPTError(errors.NewGPTError("unmarshal_response", err))
		return fmt.Sprintf("[ERROR] GPT content: %s", err)
	}

	if len(gptContent.Choices) == 0 {
		logger.GPTError(errors.NewGPTError("no_choices", fmt.Errorf("no choices in response")))
		return "[ERROR] No response from GPT"
	}

	response := gptContent.Choices[0].Message.Content
	logger.GPTResponse(response)
	return response
}

// Request is the function that can be reassigned for testing
var Request = DefaultRequest

// RequestWithContext is the function for requests with conversation context
var RequestWithContext = DefaultRequestWithContext
