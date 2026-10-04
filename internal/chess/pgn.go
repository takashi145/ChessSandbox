package chess

import (
	"errors"
	"io"

	lib "github.com/corentings/chess/v2"
)

type PGNGame struct {
	White, Black, Result, Date string
	StartFEN                   string
	Moves                      []string
}

// ParsePGNGames reads every game in a PGN. Only the main line of each game is kept.
func ParsePGNGames(r io.Reader) ([]PGNGame, error) {
	scanner := lib.NewScanner(r)

	var games []PGNGame
	for scanner.HasNext() {
		game, err := scanner.ParseNext()
		if err != nil {
			return nil, err
		}

		positions := game.Positions()
		var moves []string
		for i, move := range game.Moves() {
			moves = append(moves, lib.AlgebraicNotation{}.Encode(positions[i], move))
		}
		games = append(games, PGNGame{
			White:    game.GetTagPair("White"),
			Black:    game.GetTagPair("Black"),
			Result:   game.GetTagPair("Result"),
			Date:     game.GetTagPair("Date"),
			StartFEN: positions[0].String(),
			Moves:    moves,
		})
	}

	if len(games) == 0 {
		return nil, errors.New("no game found in PGN data")
	}
	return games, nil
}

// PGN exports the whole line, not only the moves up to the current position.
func (h *History) PGN() string {
	game := h.rebuild(len(h.moves))
	if h.startFEN != StandardFEN {
		game.AddTagPair("SetUp", "1")
		game.AddTagPair("FEN", h.startFEN)
	}
	return game.String()
}
