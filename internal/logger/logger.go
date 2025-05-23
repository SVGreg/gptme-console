package logger

import (
	"log"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	currentLevel = LevelInfo
	debugLogger  = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	infoLogger   = log.New(os.Stdout, "[INFO]  ", log.LstdFlags)
	warnLogger   = log.New(os.Stderr, "[WARN]  ", log.LstdFlags)
	errorLogger  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
)

// SetLevel sets the logging level
func SetLevel(level Level) {
	currentLevel = level
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if currentLevel <= LevelDebug {
		debugLogger.Printf(format, args...)
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if currentLevel <= LevelInfo {
		infoLogger.Printf(format, args...)
	}
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	if currentLevel <= LevelWarn {
		warnLogger.Printf(format, args...)
	}
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	if currentLevel <= LevelError {
		errorLogger.Printf(format, args...)
	}
}

// Fatal logs an error message and exits
func Fatal(format string, args ...interface{}) {
	errorLogger.Printf(format, args...)
	os.Exit(1)
}

// Session-specific logging functions
func SessionCreated(name string) {
	Info("Session created: %s", name)
}

func SessionLoaded(name string) {
	Debug("Session loaded: %s", name)
}

func SessionSaved(name string) {
	Debug("Session saved: %s", name)
}

func SessionDeleted(name string) {
	Info("Session deleted: %s", name)
}

func SessionNotFound(name string) {
	Warn("Session not found: %s", name)
}

// Config-specific logging functions
func ConfigLoaded(path string) {
	Debug("Config loaded from: %s", path)
}

func ConfigSaved(path string) {
	Info("Config saved to: %s", path)
}

func ConfigNotFound(path string) {
	Warn("Config not found at: %s", path)
}

// GPT-specific logging functions
func GPTRequest(question string) {
	Debug("GPT request: %s", question)
}

func GPTResponse(response string) {
	Debug("GPT response received (length: %d)", len(response))
}

func GPTError(err error) {
	Error("GPT error: %v", err)
}
