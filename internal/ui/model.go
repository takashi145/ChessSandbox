package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/takashi145/chess-sandbox/internal/chess"
)

const commandList = ":flip :fen :home :end :quit"

type mode int

const (
	modeChoose mode = iota
	modeFEN
	modeBoard
	modeConfirm
)

var startChoices = []string{"Continue", "New game (standard position)", "New game (from FEN)"}

// Input loop. Typed text is a SAN move; commands start with ':'. Arrow/Home/End keys navigate when the input is empty.
type Model struct {
	history *chess.History
	store   *chess.SessionStore

	mode    mode
	saved   *chess.History
	cursor  int
	input   []rune
	pending string
	message string
	flipped bool
}

func New(store *chess.SessionStore, fen string) (Model, error) {
	m := Model{store: store, mode: modeBoard}

	if fen != "" {
		h, err := chess.NewHistory(fen)
		if err != nil {
			return m, err
		}
		m.history = h
		return m, nil
	}

	if saved := store.Load(); saved != nil {
		m.saved = saved
		m.mode = modeChoose
		return m, nil
	}

	m.history = newStandardHistory()
	return m, nil
}

// Starts on `h`, ignoring any saved session.
func NewFromHistory(store *chess.SessionStore, h *chess.History) Model {
	return Model{store: store, mode: modeBoard, history: h}
}

func newStandardHistory() *chess.History {
	h, err := chess.NewHistory(chess.StandardFEN)
	if err != nil {
		panic(err)
	}
	return h
}

func (m Model) Init() tea.Cmd { return tea.ClearScreen }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if key.String() == "ctrl+c" {
		return m.quit()
	}

	switch m.mode {
	case modeChoose:
		return m.updateChoose(key)
	case modeFEN:
		return m.updateFEN(key)
	case modeConfirm:
		return m.updateConfirm(key)
	default:
		return m.updateBoard(key)
	}
}

func (m Model) updateChoose(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		m.cursor = (m.cursor + len(startChoices) - 1) % len(startChoices)
	case "down", "j":
		m.cursor = (m.cursor + 1) % len(startChoices)
	case "esc":
		return m, tea.Quit
	case "enter":
		switch m.cursor {
		case 0:
			m.history = m.saved
			m.mode = modeBoard
		case 1:
			m.history = newStandardHistory()
			m.mode = modeBoard
		default:
			m.mode = modeFEN
		}
	}
	return m, nil
}

func (m Model) updateFEN(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		return m, tea.Quit
	case "enter":
		h, err := chess.NewHistory(strings.TrimSpace(string(m.input)))
		if err != nil {
			m.message = "Invalid FEN."
			return m, nil
		}
		m.history = h
		m.input = nil
		m.message = ""
		m.mode = modeBoard
		return m, nil
	}
	m.input = editInput(m.input, key)
	return m, nil
}

func (m Model) updateConfirm(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	san := m.pending
	m.pending = ""
	m.mode = modeBoard

	if key.String() != "y" {
		m.message = "Cancelled"
		return m, nil
	}

	m.playMove(san)
	m.save()
	return m, nil
}

func (m Model) updateBoard(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		return m.quit()
	case "enter":
		return m.submit()
	}

	if len(m.input) == 0 && m.navigate(key.String()) {
		m.message = ""
		m.save()
		return m, nil
	}

	m.input = editInput(m.input, key)
	return m, nil
}

func (m Model) navigate(key string) bool {
	switch key {
	case "left":
		m.history.Back()
	case "right":
		m.history.Forward()
	case "home":
		m.history.GoTo(0)
	case "end":
		m.history.GoTo(len(m.history.Moves()))
	default:
		return false
	}
	return true
}

func (m Model) submit() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(string(m.input))
	m.input = nil

	if text == "" {
		return m, nil
	}

	m.message = ""

	if command, ok := strings.CutPrefix(text, ":"); ok {
		return m.runCommand(strings.ToLower(strings.TrimSpace(command)))
	}

	if m.history.CanGoForward() {
		m.pending = text
		m.mode = modeConfirm
		return m, nil
	}

	m.playMove(text)
	m.save()
	return m, nil
}

func (m Model) runCommand(command string) (tea.Model, tea.Cmd) {
	switch command {
	case "f", "flip":
		m.flipped = !m.flipped
	case "fen":
		m.message = "FEN: " + m.history.CurrentFEN()
	case "home":
		m.history.GoTo(0)
	case "end":
		m.history.GoTo(len(m.history.Moves()))
	case "q", "quit":
		return m.quit()
	default:
		m.message = "Commands: " + commandList
	}

	m.save()
	return m, nil
}

func (m *Model) playMove(san string) {
	if !m.history.TryPlay(san) {
		m.message = "Illegal move: " + san
	}
}

func (m *Model) save() {
	if !m.store.Save(m.history) {
		m.message = "Could not save the session"
	}
}

func (m Model) quit() (tea.Model, tea.Cmd) {
	if m.history != nil {
		m.save()
	}
	return m, tea.Quit
}

func editInput(input []rune, key tea.KeyMsg) []rune {
	switch key.Type {
	case tea.KeyBackspace:
		if len(input) > 0 {
			return input[:len(input)-1]
		}
	case tea.KeySpace:
		return append(input, ' ')
	case tea.KeyRunes:
		return append(input, key.Runes...)
	}
	return input
}

func (m Model) View() string {
	var b strings.Builder

	switch m.mode {
	case modeChoose:
		m.viewChoose(&b)
	case modeFEN:
		b.WriteString("FEN: " + string(m.input) + "_\n")
		if m.message != "" {
			b.WriteString(colored("31", m.message) + "\n")
		}
	default:
		m.viewBoard(&b)
	}

	return b.String()
}

func (m Model) viewChoose(b *strings.Builder) {
	b.WriteString("A saved session was found\n\n")
	for i, choice := range startChoices {
		if i == m.cursor {
			b.WriteString(colored("36", "> "+choice) + "\n")
		} else {
			b.WriteString("  " + choice + "\n")
		}
	}
	b.WriteString("\n" + colored("90", "↑↓ select  Enter confirm  Esc quit") + "\n")
}

func (m Model) viewBoard(b *strings.Builder) {
	s := m.history.Snapshot()

	b.WriteString(renderBoard(s, m.flipped))
	b.WriteString("\n\n")
	b.WriteString(describeLastMove(s) + "\n")
	b.WriteString(colored("90", m.hint()) + "\n")

	if m.message != "" {
		b.WriteString(colored("33", m.message) + "\n")
	}

	if m.mode == modeConfirm {
		b.WriteString(colored("31", "This will discard the moves ahead. Continue? (y/N)"))
		return
	}
	b.WriteString("> " + string(m.input) + "_")
}

func (m Model) hint() string {
	return fmt.Sprintf("%d/%d  ←→ move  Home/End  %s  Esc save and quit",
		m.history.Position(), len(m.history.Moves()), commandList)
}

func describeLastMove(s chess.Snapshot) string {
	if s.SAN == "" {
		return "(start position)"
	}

	dots := "..."
	if s.Mover == chess.White {
		dots = "."
	}

	suffix := ""
	if s.IsCheckmate {
		suffix = " checkmate"
	} else if s.IsCheck {
		suffix = " check"
	}
	return fmt.Sprintf("%d%s %s%s", s.MoveNumber, dots, s.SAN, suffix)
}

func colored(code, text string) string {
	return "\x1b[" + code + "m" + text + reset
}
