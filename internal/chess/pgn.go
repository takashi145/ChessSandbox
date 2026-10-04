package chess

import (
	"io"

	lib "github.com/corentings/chess/v2"
)

func ParsePGN(r io.Reader) (startFEN string, moves []string, err error) {
	opt, err := lib.PGN(r)
	if err != nil {
		return "", nil, err
	}
	game := lib.NewGame(opt)

	positions := game.Positions()
	for i, move := range game.Moves() {
		moves = append(moves, lib.AlgebraicNotation{}.Encode(positions[i], move))
	}
	return positions[0].String(), moves, nil
}
