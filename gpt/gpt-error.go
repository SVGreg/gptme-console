/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package gpt

type GPTError struct {
	Error GPTErrorPayload
}

type GPTErrorPayload struct {
	Message string
	Type    string
	Param   string
	Code    string
}
