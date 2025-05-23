package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	SessionDir  = ".gptme-sessions"
	SessionFile = "sessions.json"
)

type Session struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
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
	return filepath.Join(home, SessionDir)
}

func NewSessionManager() (*SessionManager, error) {
	sm := &SessionManager{}
	if err := sm.load(); err != nil {
		return sm, nil // Return empty manager if no sessions exist
	}
	return sm, nil
}

func (sm *SessionManager) load() error {
	path := filepath.Join(getSessionsDir(), SessionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, sm)
}

func (sm *SessionManager) save() error {
	sessionDir := getSessionsDir()
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(sessionDir, SessionFile)
	data, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (sm *SessionManager) sessionPath(name string) string {
	return filepath.Join(getSessionsDir(), name+".json")
}

func (sm *SessionManager) loadSession(name string) (*Session, error) {
	path := sm.sessionPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (sm *SessionManager) saveSession(session *Session) error {
	if err := os.MkdirAll(getSessionsDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.sessionPath(session.Name), data, 0644)
}

func (sm *SessionManager) CreateSession(name string) (*Session, error) {
	// Check if session with this name already exists
	if _, err := sm.loadSession(name); err == nil {
		return nil, fmt.Errorf("session '%s' already exists", name)
	}

	session := Session{
		Name:      name,
		CreatedAt: time.Now(),
		Messages:  []Message{},
	}

	if err := sm.saveSession(&session); err != nil {
		return nil, err
	}

	sm.Current = name
	if err := sm.save(); err != nil {
		return nil, err
	}
	return &session, nil
}

func (sm *SessionManager) GetSession(name string) *Session {
	session, err := sm.loadSession(name)
	if err != nil {
		return nil
	}
	return session
}

func (sm *SessionManager) AddMessage(sessionName string, role, content string) error {
	session := sm.GetSession(sessionName)
	if session == nil {
		return fmt.Errorf("session '%s' not found", sessionName)
	}
	session.Messages = append(session.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	return sm.saveSession(session)
}

func (sm *SessionManager) ListSessions() ([]Session, error) {
	sessionDir := getSessionsDir()
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == SessionFile {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		session, err := sm.loadSession(name)
		if err != nil {
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
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == SessionFile {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(sessionDir, entry.Name())
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	sm.Current = ""
	return sm.save()
}

func (sm *SessionManager) SetCurrent(name string) error {
	session := sm.GetSession(name)
	if session == nil {
		return fmt.Errorf("session '%s' not found", name)
	}
	sm.Current = name
	return sm.save()
}
