package ui

import (
	"fmt"
	"strings"

	"github.com/takashi145/chess-sandbox/internal/chess"
)

// Both sides use the same glyphs; the side is distinguished by color only.
const (
	whiteColor          = "37"
	blackColor          = "38;5;208"
	highlightBackground = "48;5;59"
	previewBackground   = "48;5;24"
	reset               = "\x1b[0m"
)

type squareMark int

const (
	markNone squareMark = iota
	markLast
	markPreview
)

var glyphs = map[chess.PieceKind]rune{
	chess.King:   '♔',
	chess.Queen:  '♕',
	chess.Rook:   '♖',
	chess.Bishop: '♗',
	chess.Knight: '♘',
	chess.Pawn:   '♙',
}

// Returns ANSI-colored text for the 8x8 board only (caller adds header/footer).
func renderBoard(s chess.Snapshot, flipped bool, mark squareMark) string {
	files := fileLabels(flipped)
	var b strings.Builder

	fmt.Fprintf(&b, "    %s\n", files)
	b.WriteString("  ┌─────────────────┐\n")

	for _, rank := range ranks(flipped) {
		fmt.Fprintf(&b, "%d │ ", rank+1)

		for _, file := range fileOrder(flipped) {
			sq := chess.Square{File: file, Rank: rank}
			squareState := markNone
			if (s.From != nil && *s.From == sq) || (s.To != nil && *s.To == sq) {
				squareState = mark
			}

			b.WriteString(renderSquare(s.Board[file][rank], squareState))
			b.WriteByte(' ')
		}

		b.WriteString("│\n")
	}

	b.WriteString("  └─────────────────┘\n")
	fmt.Fprintf(&b, "    %s", files)

	return b.String()
}

func renderSquare(piece *chess.Piece, mark squareMark) string {
	background := markBackground(mark)

	if piece == nil {
		if background != "" {
			return "\x1b[" + background + "m " + reset
		}
		return " "
	}

	color := whiteColor
	if piece.Color == chess.Black {
		color = blackColor
	}
	if background != "" {
		color += ";" + background
	}
	return fmt.Sprintf("\x1b[%sm%c%s", color, glyphs[piece.Kind], reset)
}

func markBackground(mark squareMark) string {
	switch mark {
	case markLast:
		return highlightBackground
	case markPreview:
		return previewBackground
	}
	return ""
}

func ranks(flipped bool) []int {
	if flipped {
		return []int{0, 1, 2, 3, 4, 5, 6, 7}
	}
	return []int{7, 6, 5, 4, 3, 2, 1, 0}
}

func fileOrder(flipped bool) []int {
	if flipped {
		return []int{7, 6, 5, 4, 3, 2, 1, 0}
	}
	return []int{0, 1, 2, 3, 4, 5, 6, 7}
}

func fileLabels(flipped bool) string {
	labels := make([]string, 0, 8)
	for _, f := range fileOrder(flipped) {
		labels = append(labels, string(rune('a'+f)))
	}
	return strings.Join(labels, " ")
}
