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
