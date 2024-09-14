/*
Copyright © 2024 GPTMe
*/

package gpt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SVGreg/gptme-console/config"
)

const (
	completionUrl string = "https://api.openai.com/v1/chat/completions"
	headerAuthKey string = "Authorization"
	headerAuthPrefix string = "Bearer "
	contentTypeKey string = "Content-Type"
	contentTypeValue = "application/json"
	organizationKey = "OpenAI-Organization"
	projectKey = "OpenAI-Project"
)

func Request(question string, config config.Config) {
	bodyString := `{
     "model": "gpt-4o-mini",
     "messages": [{"role": "user", "content": "%s"}],
     "temperature": 0.5
}`
    reader := strings.NewReader(fmt.Sprintf(bodyString, question))

	request, err := http.NewRequest(http.MethodPost, completionUrl, reader)
	if err != nil {
        log.Println("Error creating request:", err)
        return
    }

	request.Header.Set(contentTypeKey, contentTypeValue)
	request.Header.Set(organizationKey, config.OrganizationId)
	request.Header.Set(projectKey, config.ProjectId)
	request.Header.Set(headerAuthKey, headerAuthPrefix + config.APIKey)

	log.Println("Headers: ", request.Header)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("HTTP error: ", err)
	}

	resBody, err := io.ReadAll(res.Body)
	log.Printf("%s", resBody)
}
