/*
Copyright © 2024 GPTMe
*/

package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SVGreg/gptme-console/config"
)

const (
	completionUrl    string = "https://api.openai.com/v1/chat/completions"
	headerAuthKey    string = "Authorization"
	headerAuthPrefix string = "Bearer "
	contentTypeKey   string = "Content-Type"
	contentTypeValue        = "application/json"
	organizationKey         = "OpenAI-Organization"
	projectKey              = "OpenAI-Project"
)

func Request(question string, config config.Config) string {
	bodyString := `{
     "model": "gpt-4o-mini",
     "messages": [{"role": "user", "content": "%s"}],
     "temperature": 0.5
}`
	reader := strings.NewReader(fmt.Sprintf(bodyString, question))

	request, err := http.NewRequest(http.MethodPost, completionUrl, reader)
	if err != nil {
		log.Println("Can't create request:", err)
		return "[ERROR] Can't create request"
	}

	request.Header.Set(contentTypeKey, contentTypeValue)
	request.Header.Set(organizationKey, config.OrganizationId)
	request.Header.Set(projectKey, config.ProjectId)
	request.Header.Set(headerAuthKey, headerAuthPrefix+config.APIKey)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("HTTP error: ", err)
	}

	resBody, _ := io.ReadAll(res.Body)

	var gptError GPTError
	if _ = json.Unmarshal(resBody, &gptError); len(gptError.Error.Message) > 0 {
		return fmt.Sprintf("[ERROR] GPT reports error: %s", gptError.Error.Message)
	}

	var gptContent GPTResponse
	if err = json.Unmarshal(resBody, &gptContent); err != nil {
		return fmt.Sprintf("[ERROR] GPT content: %s", err)
	}
	return gptContent.Choices[0].Message.Content
}
