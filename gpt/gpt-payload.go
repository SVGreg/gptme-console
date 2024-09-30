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
