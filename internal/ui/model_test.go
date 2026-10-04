package ui

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/takashi145/chess-sandbox/internal/chess"
)

func newStore(t *testing.T) *chess.SessionStore {
	t.Helper()
	return chess.NewSessionStore(filepath.Join(t.TempDir(), "session.json"))
}

func newModel(t *testing.T, store *chess.SessionStore, fen string) Model {
	t.Helper()
	m, err := New(store, fen)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func press(m Model, keys ...tea.KeyMsg) Model {
	for _, k := range keys {
		next, _ := m.Update(k)
		m = next.(Model)
	}
	return m
}

func typed(m Model, text string) Model {
	return press(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)}, tea.KeyMsg{Type: tea.KeyEnter})
}

var (
	left  = tea.KeyMsg{Type: tea.KeyLeft}
	right = tea.KeyMsg{Type: tea.KeyRight}
	home  = tea.KeyMsg{Type: tea.KeyHome}
	end   = tea.KeyMsg{Type: tea.KeyEnd}
	down  = tea.KeyMsg{Type: tea.KeyDown}
	enter = tea.KeyMsg{Type: tea.KeyEnter}
	esc   = tea.KeyMsg{Type: tea.KeyEsc}
)

func runes(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func TestPlayingMovesAndNavigating(t *testing.T) {
	m := newModel(t, newStore(t), "")
	m = typed(m, "e4")
	m = typed(m, "e5")

	if got := m.history.Moves(); !slices.Equal(got, []string{"e4", "e5"}) {
		t.Fatalf("moves = %v", got)
	}

	m = press(m, left)
	if m.history.Position() != 1 {
		t.Errorf("after left: pos=%d", m.history.Position())
	}
	m = press(m, home)
	if m.history.Position() != 0 {
		t.Errorf("after home: pos=%d", m.history.Position())
	}
	m = press(m, end)
	if m.history.Position() != 2 {
		t.Errorf("after end: pos=%d", m.history.Position())
	}
}

func TestArrowKeysDoNotNavigateWhileTyping(t *testing.T) {
	m := newModel(t, newStore(t), "")
	m = typed(m, "e4")

	m = press(m, runes("N"), left)

	if m.history.Position() != 1 {
		t.Errorf("pos=%d, want 1", m.history.Position())
	}
	if string(m.input) != "N" {
		t.Errorf("input = %q, want N", string(m.input))
	}
}

func TestIllegalMoveShowsMessageAndKeepsPosition(t *testing.T) {
	m := newModel(t, newStore(t), "")
	m = typed(m, "e5")

	if m.message != "Illegal move: e5" || m.history.Position() != 0 {
		t.Errorf("message=%q pos=%d", m.message, m.history.Position())
	}
}

func TestPlayingAfterGoingBackAsksBeforeDiscarding(t *testing.T) {
	m := newModel(t, newStore(t), "")
	m = typed(m, "e4")
	m = typed(m, "e5")
	m = press(m, left)

	m = typed(m, "c5")
	if m.mode != modeConfirm || len(m.history.Moves()) != 2 {
		t.Fatalf("mode=%v moves=%v", m.mode, m.history.Moves())
	}
	if !strings.Contains(m.View(), "discard the moves ahead") {
		t.Error("confirmation prompt not shown")
	}

	m = press(m, runes("n"))
	if m.mode != modeBoard || m.message != "Cancelled" || !slices.Equal(m.history.Moves(), []string{"e4", "e5"}) {
		t.Errorf("after cancel: mode=%v message=%q moves=%v", m.mode, m.message, m.history.Moves())
	}

	m = typed(m, "c5")
	m = press(m, runes("y"))
	if !slices.Equal(m.history.Moves(), []string{"e4", "c5"}) {
		t.Errorf("after confirm: moves=%v", m.history.Moves())
	}
}

func TestCommands(t *testing.T) {
	m := newModel(t, newStore(t), "")
	m = typed(m, "e4")

	m = typed(m, ":flip")
	if !m.flipped {
		t.Error(":flip did not flip")
	}
	m = typed(m, ":f")
	if m.flipped {
		t.Error(":f did not flip back")
	}

	m = typed(m, ":fen")
	if !strings.HasPrefix(m.message, "FEN: rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b") {
		t.Errorf("message = %q", m.message)
	}

	m = typed(m, ":home")
	if m.history.Position() != 0 {
		t.Errorf(":home: pos=%d", m.history.Position())
	}
	m = typed(m, ":end")
	if m.history.Position() != 1 {
		t.Errorf(":end: pos=%d", m.history.Position())
	}

	m = typed(m, ":bogus")
	if !strings.HasPrefix(m.message, "Commands:") {
		t.Errorf("message = %q", m.message)
	}
}

func TestQuitSavesAndEnds(t *testing.T) {
	for _, command := range []string{"", ":quit", ":q"} {
		store := newStore(t)
		m := typed(newModel(t, store, ""), "e4")

		var cmd tea.Cmd
		if command == "" {
			_, cmd = m.Update(esc)
		} else {
			m.input = []rune(command)
			_, cmd = m.Update(enter)
		}

		if cmd == nil || cmd() != tea.Quit() {
			t.Errorf("%q: expected tea.Quit", command)
		}
		if saved := store.Load(); saved == nil || !slices.Equal(saved.Moves(), []string{"e4"}) {
			t.Errorf("%q: saved = %+v", command, saved)
		}
	}
}

func TestFENFlagStartsFromThatPosition(t *testing.T) {
	store := newStore(t)
	store.Save(mustHistory(t, "e4"))
	const fen = "7k/P7/8/8/8/8/8/K7 w - - 0 1"

	m := newModel(t, store, fen)

	if m.mode != modeBoard || m.history.StartFEN() != fen {
		t.Errorf("mode=%v start=%q", m.mode, m.history.StartFEN())
	}
}

func TestInvalidFENFlagIsAnError(t *testing.T) {
	if _, err := New(newStore(t), "bad"); err == nil {
		t.Error("expected an error")
	}
}

func TestNoSavedSessionStartsNewGame(t *testing.T) {
	m := newModel(t, newStore(t), "")

	if m.mode != modeBoard || m.history.StartFEN() != chess.StandardFEN {
		t.Errorf("mode=%v start=%q", m.mode, m.history.StartFEN())
	}
}

func TestSavedSessionOffersChoices(t *testing.T) {
	store := newStore(t)
	store.Save(mustHistory(t, "e4", "e5"))

	m := newModel(t, store, "")
	if m.mode != modeChoose {
		t.Fatalf("mode = %v, want modeChoose", m.mode)
	}

	cont := press(m, enter)
	if cont.mode != modeBoard || len(cont.history.Moves()) != 2 {
		t.Errorf("continue: mode=%v moves=%v", cont.mode, cont.history.Moves())
	}

	fresh := press(m, down, enter)
	if fresh.mode != modeBoard || len(fresh.history.Moves()) != 0 {
		t.Errorf("new game: mode=%v moves=%v", fresh.mode, fresh.history.Moves())
	}

	custom := press(m, down, down, enter)
	if custom.mode != modeFEN {
		t.Fatalf("mode = %v, want modeFEN", custom.mode)
	}
	custom = typed(custom, "bad")
	if custom.mode != modeFEN || custom.message != "Invalid FEN." {
		t.Errorf("invalid fen: mode=%v message=%q", custom.mode, custom.message)
	}
	custom = press(custom, tea.KeyMsg{Type: tea.KeyBackspace}, tea.KeyMsg{Type: tea.KeyBackspace}, tea.KeyMsg{Type: tea.KeyBackspace})
	custom = typed(custom, "7k/P7/8/8/8/8/8/K7 w - - 0 1")
	if custom.mode != modeBoard || custom.history.StartFEN() != "7k/P7/8/8/8/8/8/K7 w - - 0 1" {
		t.Errorf("valid fen: mode=%v start=%q", custom.mode, custom.history.StartFEN())
	}
}

func TestChoosingDoesNotOverwriteTheSavedSession(t *testing.T) {
	store := newStore(t)
	store.Save(mustHistory(t, "e4", "e5"))

	press(newModel(t, store, ""), down, enter)

	if saved := store.Load(); saved == nil || len(saved.Moves()) != 2 {
		t.Errorf("saved session was overwritten: %+v", saved)
	}
}

func TestEscInChooseQuitsAndKeepsTheSavedSession(t *testing.T) {
	store := newStore(t)
	store.Save(mustHistory(t, "e4", "e5"))
	m := newModel(t, store, "")

	_, cmd := m.Update(esc)

	if cmd == nil || cmd() != tea.Quit() {
		t.Error("expected tea.Quit")
	}
	if saved := store.Load(); saved == nil || len(saved.Moves()) != 2 {
		t.Errorf("saved session changed: %+v", saved)
	}
}

func mustHistory(t *testing.T, moves ...string) *chess.History {
	t.Helper()
	h, err := chess.NewHistory(chess.StandardFEN)
	if err != nil {
		t.Fatal(err)
	}
	for _, mv := range moves {
		if !h.TryPlay(mv) {
			t.Fatalf("expected %s to be legal", mv)
		}
	}
	return h
}
