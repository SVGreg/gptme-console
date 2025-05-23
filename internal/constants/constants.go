package constants

const (
	// Application
	AppName = "gptme-console"

	// Configuration
	ConfigFileName = ".gptme-config.json"

	// Sessions
	SessionDir     = ".gptme-sessions"
	SessionFile    = "sessions.json"
	SessionFileExt = ".json"

	// Limits
	MaxQuestionWords = 30

	// GPT API
	GPTModel       = "gpt-4o-mini"
	GPTTemperature = 0.5
	CompletionURL  = "https://api.openai.com/v1/chat/completions"

	// HTTP Headers
	HeaderAuth         = "Authorization"
	HeaderAuthPrefix   = "Bearer "
	HeaderContentType  = "Content-Type"
	HeaderOrganization = "OpenAI-Organization"
	HeaderProject      = "OpenAI-Project"
	ContentTypeJSON    = "application/json"

	// Environment Variables
	EnvAPIKey     = "GPTME_API_KEY"
	EnvOrgID      = "GPTME_ORG_ID"
	EnvProjectID  = "GPTME_PROJECT_ID"
	EnvConfigPath = "GPTME_CONFIG_PATH"

	// File Permissions
	DirPerms  = 0755
	FilePerms = 0644

	// Formatting
	MarkdownWidth  = 120
	MarkdownIndent = 2
)
