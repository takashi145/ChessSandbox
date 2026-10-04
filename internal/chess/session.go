package chess

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type sessionData struct {
	StartFen string
	Moves    []string
	Position int
}

type SessionStore struct {
	path string
}

func NewSessionStore(path string) *SessionStore { return &SessionStore{path: path} }

func DefaultSessionPath() string {
	return filepath.Join(localDataDir(), "ChessSandbox", "session.json")
}

func DefaultSessionStore() *SessionStore { return NewSessionStore(DefaultSessionPath()) }

func localDataDir() string {
	home, _ := os.UserHomeDir()

	if runtime.GOOS == "windows" {
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return dir
		}
		return filepath.Join(home, "AppData", "Local")
	}

	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(home, ".local", "share")
}

// Returns false if the session could not be written. The previous file is left intact in that case.
func (s *SessionStore) Save(h *History) bool {
	content, err := json.Marshal(sessionData{StartFen: h.startFEN, Moves: h.Moves(), Position: h.position})
	if err != nil {
		return false
	}

	temp := s.path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return false
	}
	if err := os.WriteFile(temp, content, 0o644); err != nil {
		return false
	}
	if err := os.Rename(temp, s.path); err != nil {
		os.Remove(temp)
		return false
	}
	return true
}

// Returns nil if the file is missing, unreadable or does not describe a legal game.
func (s *SessionStore) Load() *History {
	content, err := os.ReadFile(s.path)
	if err != nil {
		return nil
	}

	var data sessionData
	if err := json.Unmarshal(content, &data); err != nil || data.Moves == nil {
		return nil
	}

	h, err := RestoreHistory(data.StartFen, data.Moves, data.Position)
	if err != nil {
		return nil
	}
	return h
}
