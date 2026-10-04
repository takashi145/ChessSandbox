package chess

import (
	"errors"
	"strings"

	lib "github.com/corentings/chess/v2"
)

const StandardFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type Side int

const (
	White Side = iota
	Black
)

type PieceKind int

const (
	Pawn PieceKind = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece struct {
	Color Side
	Kind  PieceKind
}

type Square struct {
	File int
	Rank int
}

type Snapshot struct {
	Board       [8][8]*Piece
	MoveNumber  int
	SideToMove  Side
	Mover       Side
	SAN         string
	From, To    *Square
	IsCheck     bool
	IsCheckmate bool
}

func newGame(fen string) (*lib.Game, error) {
	opt, err := lib.FEN(strings.TrimSpace(fen))
	if err != nil {
		return nil, err
	}
	game := lib.NewGame(opt)
	if err := checkPieces(game.Position()); err != nil {
		return nil, err
	}
	return game, nil
}

func checkPieces(pos *lib.Position) error {
	var kings [2]int
	board := pos.Board()

	for file := range 8 {
		for rank := range 8 {
			p := board.Piece(lib.NewSquare(lib.File(file), lib.Rank(rank)))
			if p == lib.NoPiece {
				continue
			}
			if p.Type() == lib.King {
				kings[toSide(p.Color())]++
			}
			if p.Type() == lib.Pawn && (rank == 0 || rank == 7) {
				return errors.New("pawn on the first or last rank")
			}
		}
	}

	if kings != [2]int{1, 1} {
		return errors.New("each side needs exactly one king")
	}
	return nil
}

func playSAN(game *lib.Game, san string) (string, bool) {
	san = strings.TrimSpace(san)
	if san == "" {
		return "", false
	}
	if err := game.PushNotationMove(san, lib.AlgebraicNotation{}, nil); err != nil {
		return "", false
	}
	moves := game.Moves()
	positions := game.Positions()
	last := moves[len(moves)-1]
	return lib.AlgebraicNotation{}.Encode(positions[len(positions)-2], last), true
}

func currentFEN(game *lib.Game) string {
	return game.Position().String()
}

// Snapshot of `game` as it stands, after its last executed move (or the start position if there is none).
func snapshot(game *lib.Game) Snapshot {
	pos := game.Position()
	s := Snapshot{
		Board:      captureBoard(pos),
		SideToMove: toSide(pos.Turn()),
	}

	moves := game.Moves()
	if len(moves) == 0 {
		return s
	}

	last := moves[len(moves)-1]
	positions := game.Positions()
	before := positions[len(positions)-2]
	s.Mover = toSide(before.Turn())
	s.MoveNumber = (before.Ply() + 1) / 2
	s.SAN = lib.AlgebraicNotation{}.Encode(before, last)
	s.From = toSquare(last.S1())
	s.To = toSquare(last.S2())
	s.IsCheck = last.HasTag(lib.Check)
	s.IsCheckmate = game.Method() == lib.Checkmate
	return s
}

func captureBoard(pos *lib.Position) [8][8]*Piece {
	var grid [8][8]*Piece
	board := pos.Board()
	for file := range 8 {
		for rank := range 8 {
			p := board.Piece(lib.NewSquare(lib.File(file), lib.Rank(rank)))
			if p == lib.NoPiece {
				continue
			}
			grid[file][rank] = &Piece{Color: toSide(p.Color()), Kind: toKind(p.Type())}
		}
	}
	return grid
}

func toSquare(sq lib.Square) *Square {
	return &Square{File: int(sq.File()), Rank: int(sq.Rank())}
}

func toSide(c lib.Color) Side {
	if c == lib.White {
		return White
	}
	return Black
}

func toKind(t lib.PieceType) PieceKind {
	switch t {
	case lib.Knight:
		return Knight
	case lib.Bishop:
		return Bishop
	case lib.Rook:
		return Rook
	case lib.Queen:
		return Queen
	case lib.King:
		return King
	default:
		return Pawn
	}
}
