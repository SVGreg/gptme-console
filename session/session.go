package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	SessionDir  = ".gptme_sessions"
	SessionFile = "sessions.json"
)

type Session struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type SessionManager struct {
	Sessions []Session `json:"sessions"`
	Current  int       `json:"current"`
}

func NewSessionManager() (*SessionManager, error) {
	sm := &SessionManager{}
	if err := sm.load(); err != nil {
		return sm, nil // Return empty manager if no sessions exist
	}
	return sm, nil
}

func (sm *SessionManager) load() error {
	path := filepath.Join(SessionDir, SessionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, sm)
}

func (sm *SessionManager) save() error {
	if err := os.MkdirAll(SessionDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(SessionDir, SessionFile)
	data, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (sm *SessionManager) CreateSession() *Session {
	session := Session{
		ID:        len(sm.Sessions) + 1,
		CreatedAt: time.Now(),
		Messages:  []Message{},
	}
	sm.Sessions = append(sm.Sessions, session)
	sm.Current = session.ID
	sm.save()
	return &session
}

func (sm *SessionManager) GetSession(id int) *Session {
	for i := range sm.Sessions {
		if sm.Sessions[i].ID == id {
			return &sm.Sessions[i]
		}
	}
	return nil
}

func (sm *SessionManager) AddMessage(sessionID int, role, content string) error {
	session := sm.GetSession(sessionID)
	if session == nil {
		return fmt.Errorf("session %d not found", sessionID)
	}
	session.Messages = append(session.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	return sm.save()
}

func (sm *SessionManager) ListSessions() []Session {
	return sm.Sessions
}

func (sm *SessionManager) Clean() error {
	sm.Sessions = []Session{}
	sm.Current = 0
	return sm.save()
}
