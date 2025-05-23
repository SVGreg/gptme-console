/*
Copyright (c) 2024 SVGreg <git@svgreg.net>

This software is licensed under the MIT License.
See the LICENSE file for details.
*/
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SVGreg/gptme-console/internal/constants"
	"github.com/SVGreg/gptme-console/internal/errors"
	"github.com/SVGreg/gptme-console/internal/logger"
	"github.com/SVGreg/gptme-console/internal/validation"
)

type Session struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type SessionManager struct {
	Current string `json:"current"`
}

func getSessionsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, constants.SessionDir)
}

func NewSessionManager() (*SessionManager, error) {
	sm := &SessionManager{}
	if err := sm.load(); err != nil {
		logger.Debug("No existing session manager found, creating new one")
		return sm, nil // Return empty manager if no sessions exist
	}
	return sm, nil
}

func (sm *SessionManager) load() error {
	path := filepath.Join(getSessionsDir(), constants.SessionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, sm)
}

func (sm *SessionManager) save() error {
	sessionDir := getSessionsDir()
	if err := os.MkdirAll(sessionDir, constants.DirPerms); err != nil {
		return errors.NewSessionError("", "mkdir", err)
	}
	path := filepath.Join(sessionDir, constants.SessionFile)
	data, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		return errors.NewSessionError("", "marshal", err)
	}
	return os.WriteFile(path, data, constants.FilePerms)
}

func (sm *SessionManager) sessionPath(name string) string {
	return filepath.Join(getSessionsDir(), name+constants.SessionFileExt)
}

func (sm *SessionManager) loadSession(name string) (*Session, error) {
	path := sm.sessionPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewSessionError(name, "read", err)
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.NewSessionError(name, "unmarshal", err)
	}
	logger.SessionLoaded(name)
	return &session, nil
}

func (sm *SessionManager) saveSession(session *Session) error {
	if err := os.MkdirAll(getSessionsDir(), constants.DirPerms); err != nil {
		return errors.NewSessionError(session.Name, "mkdir", err)
	}

	session.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return errors.NewSessionError(session.Name, "marshal", err)
	}

	err = os.WriteFile(sm.sessionPath(session.Name), data, constants.FilePerms)
	if err != nil {
		return errors.NewSessionError(session.Name, "write", err)
	}

	logger.SessionSaved(session.Name)
	return nil
}

func (sm *SessionManager) CreateSession(name string, description string, tags []string) (*Session, error) {
	// Validate session name
	if err := validation.ValidateSessionName(name); err != nil {
		return nil, err
	}

	// Check if session with this name already exists
	if _, err := sm.loadSession(name); err == nil {
		return nil, errors.NewSessionError(name, "create",
			fmt.Errorf("session already exists"))
	}

	now := time.Now()
	session := Session{
		Name:        name,
		Description: description,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
		Messages:    []Message{},
	}

	if err := sm.saveSession(&session); err != nil {
		return nil, err
	}

	sm.Current = name
	if err := sm.save(); err != nil {
		return nil, errors.NewSessionError(name, "set_current", err)
	}

	logger.SessionCreated(name)
	return &session, nil
}

func (sm *SessionManager) GetSession(name string) *Session {
	session, err := sm.loadSession(name)
	if err != nil {
		logger.SessionNotFound(name)
		return nil
	}
	return session
}

func (sm *SessionManager) AddMessage(sessionName string, role, content string) error {
	session := sm.GetSession(sessionName)
	if session == nil {
		return errors.NewSessionError(sessionName, "add_message",
			fmt.Errorf("session not found"))
	}

	session.Messages = append(session.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	return sm.saveSession(session)
}

func (sm *SessionManager) UpdateSessionMetadata(name, description string, tags []string) error {
	session := sm.GetSession(name)
	if session == nil {
		return errors.NewSessionError(name, "update_metadata",
			fmt.Errorf("session not found"))
	}

	session.Description = description
	session.Tags = tags
	return sm.saveSession(session)
}

func (sm *SessionManager) ListSessions() ([]Session, error) {
	sessionDir := getSessionsDir()
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, errors.NewSessionError("", "list", err)
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == constants.SessionFile {
			continue
		}
		if filepath.Ext(entry.Name()) != constants.SessionFileExt {
			continue
		}

		name := entry.Name()[:len(entry.Name())-len(constants.SessionFileExt)]
		session, err := sm.loadSession(name)
		if err != nil {
			logger.Warn("Failed to load session %s: %v", name, err)
			continue
		}
		sessions = append(sessions, *session)
	}
	return sessions, nil
}

func (sm *SessionManager) Clean() error {
	sessionDir := getSessionsDir()
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return errors.NewSessionError("", "clean", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == constants.SessionFile {
			continue
		}
		if filepath.Ext(entry.Name()) != constants.SessionFileExt {
			continue
		}

		name := entry.Name()[:len(entry.Name())-len(constants.SessionFileExt)]
		path := filepath.Join(sessionDir, entry.Name())
		if err := os.Remove(path); err != nil {
			return errors.NewSessionError(name, "delete", err)
		}
		logger.SessionDeleted(name)
	}

	sm.Current = ""
	return sm.save()
}

func (sm *SessionManager) SetCurrent(name string) error {
	session := sm.GetSession(name)
	if session == nil {
		return errors.NewSessionError(name, "set_current",
			fmt.Errorf("session not found"))
	}
	sm.Current = name
	return sm.save()
}

// GetConversationHistory returns the message history for GPT context
func (sm *SessionManager) GetConversationHistory(sessionName string) ([]Message, error) {
	session := sm.GetSession(sessionName)
	if session == nil {
		return nil, errors.NewSessionError(sessionName, "get_history",
			fmt.Errorf("session not found"))
	}
	return session.Messages, nil
}
