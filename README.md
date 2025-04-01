# GPTme Console

A command-line interface for communicating with GPT models, featuring session management for persistent conversations.

## Features

- Interactive chat with GPT models
- Session management for persistent conversations
- Markdown rendering for responses
- Configurable API settings
- Session history viewing and management

## Installation

```bash
go install github.com/SVGreg/gptme-console@latest
```

## Configuration

Before using the tool, you need to set up your OpenAI API key and other parameters. It is possible to specify `path` to configuration file using `-p` parameter. 

### Command

Use `init` command to initialize the project
```bash
gptme-console init ...
```

### Manually
 Create a configuration file at `~/.gptme-config.json` with the following structure:

```json
{
    "OrganizationId":"org-YYYYY",
    "ProjectId":"proj_XXXX",
    "APIKey":"sk-proj-XXXXXXX"
}
```

## Commands

### Basic Usage

```bash
gptme-console ask "Your question here"
```

This command sends a question to GPT and displays the response in markdown format. This is ad-hoc question&answer and will not be stored in session.

### Session Management

The tool provides several commands for managing conversation sessions:

#### Start a New Session

```bash
gptme-console session --start "session-name"
```

Creates a new session with the specified name and makes it the current session.

#### Use Existing Session

```bash
gptme-console session --use "session-name"
```

Switches to an existing session. Raises an error if the session doesn't exist.

#### List Sessions

```bash
gptme-console session --list
```

Displays all available sessions with their creation dates and message counts.

#### View Session History

```bash
gptme-console session --cat
```

Shows the history of the current session. Recommended to use with `less` or `more` for better readability.

#### Ask Questions in Session

```bash
gptme-console session --ask "Your question here"
```

Sends a question in the current session context. The question and response will be stored in the session history.

#### Clean Sessions

```bash
gptme-console session --clean
```

Removes all stored sessions.

## Session Storage

Sessions are stored in the `~/.gptme-sessions/` directory:
- Each session is stored in its own file: `<session-name>.json`
- The current session is tracked in `sessions.json`

## Limitations

- A current session must be selected (using `--start` or `--use`) before asking session questions
- Session names must be unique

## License

[Your License Here] 