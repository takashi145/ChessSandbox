package chess

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func sessionPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "nested", "session.json")
}

func TestSaveThenLoadRestoresSamePosition(t *testing.T) {
	h := play(t, "e4", "e5", "Nf3", "Nc6")
	h.GoTo(2)

	store := NewSessionStore(sessionPath(t))
	if !store.Save(h) {
		t.Fatal("Save failed")
	}
	loaded := store.Load()

	if loaded == nil {
		t.Fatal("Load returned nil")
	}
	if loaded.StartFEN() != h.StartFEN() || !slices.Equal(loaded.Moves(), h.Moves()) ||
		loaded.Position() != h.Position() || loaded.CurrentFEN() != h.CurrentFEN() {
		t.Errorf("loaded session differs: %v at %d", loaded.Moves(), loaded.Position())
	}
}

func TestSaveThenLoadKeepsCustomStartFEN(t *testing.T) {
	const fen = "7k/P7/8/8/8/8/8/K7 w - - 0 1"
	h, err := NewHistory(fen)
	if err != nil || !h.TryPlay("a8=Q") {
		t.Fatal("setup failed")
	}

	store := NewSessionStore(sessionPath(t))
	store.Save(h)
	loaded := store.Load()

	if loaded == nil || loaded.StartFEN() != fen || !slices.Equal(loaded.Moves(), h.Moves()) {
		t.Errorf("loaded = %+v", loaded)
	}
}

func TestSaveCreatesMissingDirectory(t *testing.T) {
	path := sessionPath(t)

	if !NewSessionStore(path).Save(play(t)) {
		t.Fatal("Save failed")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("session file missing: %v", err)
	}
}

func TestSaveOverwritesPreviousSession(t *testing.T) {
	store := NewSessionStore(sessionPath(t))
	store.Save(play(t, "e4"))

	store.Save(play(t))

	if loaded := store.Load(); loaded == nil || len(loaded.Moves()) != 0 {
		t.Errorf("loaded = %+v", loaded)
	}
}

func TestLoadMissingFileReturnsNil(t *testing.T) {
	if NewSessionStore(sessionPath(t)).Load() != nil {
		t.Error("expected nil")
	}
}

func TestLoadCorruptFileReturnsNil(t *testing.T) {
	cases := []string{
		"",
		"not json",
		"{",
		"null",
		"{}",
		`{"StartFen":"bad","Moves":[],"Position":0}`,
		`{"StartFen":"` + StandardFEN + `","Moves":["e4","Ke5"],"Position":2}`,
		`{"StartFen":"` + StandardFEN + `","Moves":["e4"],"Position":5}`,
	}

	for _, content := range cases {
		path := sessionPath(t)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		if NewSessionStore(path).Load() != nil {
			t.Errorf("Load should return nil for %q", content)
		}
	}
}

func TestLoadReadsSavedFileFormat(t *testing.T) {
	const content = `{"StartFen":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1","Moves":["e4","e5","Nf3"],"Position":2}`
	path := sessionPath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	loaded := NewSessionStore(path).Load()

	if loaded == nil || len(loaded.Moves()) != 3 || loaded.Position() != 2 {
		t.Errorf("loaded = %+v", loaded)
	}
}
