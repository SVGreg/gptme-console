/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package gpt

type GPTResponse struct {
	Id      string
	Object  string
	Created int
	Model   string
	Choices []GPTChoice
}

type GPTChoice struct {
	Index   int
	Message GPTMessage
}

type GPTMessage struct {
	Role    string
	Content string
}
