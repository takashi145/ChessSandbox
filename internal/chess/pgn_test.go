package chess

import (
	"slices"
	"strings"
	"testing"
)

func TestParsePGNReadsHeadersAndMoves(t *testing.T) {
	const pgn = `[Event "Casual Game"]
[White "Alice"]
[Black "Bob"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 1-0
`
	fen, moves, err := ParsePGN(strings.NewReader(pgn))

	if err != nil {
		t.Fatal(err)
	}
	if fen != StandardFEN {
		t.Errorf("fen = %q, want %q", fen, StandardFEN)
	}
	want := []string{"e4", "e5", "Nf3", "Nc6", "Bb5", "a6"}
	if !slices.Equal(moves, want) {
		t.Errorf("moves = %v, want %v", moves, want)
	}
}

func TestParsePGNUsesTheFENHeaderAsStartPosition(t *testing.T) {
	const startFEN = "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"
	const pgn = `[SetUp "1"]
[FEN "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"]

3... a6 4. Ba4 Nf6 *
`
	fen, moves, err := ParsePGN(strings.NewReader(pgn))

	if err != nil {
		t.Fatal(err)
	}
	if fen != startFEN {
		t.Errorf("fen = %q, want %q", fen, startFEN)
	}
	want := []string{"a6", "Ba4", "Nf6"}
	if !slices.Equal(moves, want) {
		t.Errorf("moves = %v, want %v", moves, want)
	}
}

func TestParsePGNResultCanBeRestoredAsHistory(t *testing.T) {
	const pgn = `[Event "Opera Game"]

1. e4 e5 2. Nf3 d6 3. d4 Bg4 4. dxe5 Bxf3 5. Qxf3 dxe5 6. Bc4 Nf6 7. Qb3 Qe7
8. Nc3 c6 9. Bg5 b5 10. Nxb5 cxb5 11. Bxb5+ Nbd7 12. O-O-O Rd8 13. Rxd7 Rxd7
14. Rd1 Qe6 15. Bxd7+ Nxd7 16. Qb8+ Nxb8 17. Rd8# 1-0
`
	fen, moves, err := ParsePGN(strings.NewReader(pgn))
	if err != nil {
		t.Fatal(err)
	}

	h, err := RestoreHistory(fen, moves, 0)
	if err != nil {
		t.Fatal(err)
	}

	if got := len(h.Moves()); got != 33 {
		t.Errorf("restored %d moves, want 33", got)
	}
	if !slices.Equal(h.Moves(), moves) {
		t.Errorf("restored moves differ from the parsed moves:\n got  %v\n want %v", h.Moves(), moves)
	}
}
