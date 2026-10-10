package chess

import (
	"fmt"
	"strings"

	lib "github.com/corentings/chess/v2"
)

// One line of play: start position + SAN moves + current position.
type History struct {
	startFEN string
	moves    []string
	position int
	game     *lib.Game
}

func NewHistory(startFEN string) (*History, error) {
	game, err := newGame(startFEN)
	if err != nil {
		return nil, err
	}
	return &History{startFEN: startFEN, game: game}, nil
}

func RestoreHistory(startFEN string, moves []string, position int) (*History, error) {
	h, err := NewHistory(startFEN)
	if err != nil {
		return nil, err
	}

	for _, san := range moves {
		if !h.playOn(h.game, san) {
			return nil, fmt.Errorf("illegal move %q", san)
		}
	}

	if position < 0 || position > len(h.moves) {
		return nil, fmt.Errorf("position %d out of range", position)
	}

	h.GoTo(position)
	return h, nil
}

func (h *History) StartFEN() string { return h.startFEN }
func (h *History) Moves() []string  { return append([]string{}, h.moves...) }
func (h *History) Position() int    { return h.position }
func (h *History) CanGoBack() bool  { return h.position > 0 }

// Playing a move while this is true discards the moves ahead of the current position.
func (h *History) CanGoForward() bool { return h.position < len(h.moves) }

func (h *History) CurrentFEN() string { return currentFEN(h.game) }

func (h *History) Snapshot() Snapshot { return snapshot(h.game) }

type PreviewState int

const (
	PreviewNone PreviewState = iota
	PreviewLegal
	PreviewPartial
	PreviewInvalid
)

type Preview struct {
	State    PreviewState
	From, To *Square
}

// Preview matches what TryPlay accepts, because both use the same notation decoder.
func (h *History) Preview(input string) Preview {
	input = strings.TrimSpace(input)
	if input == "" {
		return Preview{}
	}

	pos := h.game.Position()
	move, err := lib.AlgebraicNotation{}.Decode(pos, input)
	if err == nil {
		return Preview{State: PreviewLegal, From: toSquare(move.S1()), To: toSquare(move.S2())}
	}

	for _, m := range pos.ValidMoves() {
		san := strings.TrimRight(lib.AlgebraicNotation{}.Encode(pos, &m), "+#")
		coordinates := m.S1().String() + m.S2().String()
		if strings.HasPrefix(san, input) || strings.HasPrefix(coordinates, input) {
			return Preview{State: PreviewPartial}
		}
	}
	return Preview{State: PreviewInvalid}
}

// Plays a move at the current position, dropping any moves ahead. State is unchanged on failure.
func (h *History) TryPlay(san string) bool {
	probe := h.rebuild(h.position)
	if !h.playOn(probe, san) {
		return false
	}

	h.game = probe
	return true
}

func (h *History) Back()    { h.GoTo(h.position - 1) }
func (h *History) Forward() { h.GoTo(h.position + 1) }

func (h *History) GoTo(position int) {
	position = max(0, min(position, len(h.moves)))
	if position == h.position {
		return
	}

	h.game = h.rebuild(position)
	h.position = position
}

// Plays `san` on `game`, which must be at the current position, and records it.
func (h *History) playOn(game *lib.Game, san string) bool {
	played, ok := playSAN(game, san)
	if !ok {
		return false
	}

	h.moves = append(h.moves[:h.position], played)
	h.position++
	return true
}

func (h *History) rebuild(position int) *lib.Game {
	game, err := newGame(h.startFEN)
	if err != nil {
		panic(err)
	}

	for _, san := range h.moves[:position] {
		if _, ok := playSAN(game, san); !ok {
			panic(fmt.Sprintf("recorded move %q is no longer legal", san))
		}
	}

	return game
}
