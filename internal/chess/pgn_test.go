package chess

import (
	"reflect"
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
	games, err := ParsePGNGames(strings.NewReader(pgn))

	if err != nil {
		t.Fatal(err)
	}
	fen, moves := games[0].StartFEN, games[0].Moves
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
	games, err := ParsePGNGames(strings.NewReader(pgn))

	if err != nil {
		t.Fatal(err)
	}
	fen, moves := games[0].StartFEN, games[0].Moves
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
	games, err := ParsePGNGames(strings.NewReader(pgn))
	if err != nil {
		t.Fatal(err)
	}
	fen, moves := games[0].StartFEN, games[0].Moves

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

func TestPGNOfAGameStartedFromAFENUsesTheFENsMoveNumber(t *testing.T) {
	const fen = "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 23"
	h, err := NewHistory(fen)
	if err != nil {
		t.Fatal(err)
	}
	h.TryPlay("a6")
	h.TryPlay("Ba4")

	got := h.PGN()

	for _, want := range []string{`[SetUp "1"]`, `[FEN "` + fen + `"]`, "23... a6 24. Ba4"} {
		if !strings.Contains(got, want) {
			t.Errorf("PGN is missing %q:\n%s", want, got)
		}
	}
}

func TestPGNCanBeParsedBack(t *testing.T) {
	h := play(t, "e4", "e5", "Nf3", "Nc6", "Bb5")

	games, err := ParsePGNGames(strings.NewReader(h.PGN()))
	if err != nil {
		t.Fatal(err)
	}
	fen, moves := games[0].StartFEN, games[0].Moves

	if fen != h.StartFEN() {
		t.Errorf("fen = %q, want %q", fen, h.StartFEN())
	}
	if !slices.Equal(moves, h.Moves()) {
		t.Errorf("moves = %v, want %v", moves, h.Moves())
	}
}

func TestParsePGNGamesReadsEveryGame(t *testing.T) {
	const pgn = `[Event "One"]
[White "Alice"]
[Black "Bob"]
[Result "1-0"]
[Date "2026.10.01"]

1. e4 e5 2. Nf3 1-0

[Event "Two"]
[White "Carol"]
[Black "Dave"]
[Result "1/2-1/2"]

1. d4 d5 1/2-1/2

[Event "Three"]
[SetUp "1"]
[FEN "7k/P7/8/8/8/8/8/K7 w - - 0 1"]

1. a8=Q+ *
`
	games, err := ParsePGNGames(strings.NewReader(pgn))
	if err != nil {
		t.Fatal(err)
	}

	want := []PGNGame{
		{"Alice", "Bob", "1-0", "2026.10.01", StandardFEN, []string{"e4", "e5", "Nf3"}},
		{"Carol", "Dave", "1/2-1/2", "", StandardFEN, []string{"d4", "d5"}},
		{"", "", "", "", "7k/P7/8/8/8/8/8/K7 w - - 0 1", []string{"a8=Q+"}},
	}
	if !reflect.DeepEqual(games, want) {
		t.Errorf("games = %+v, want %+v", games, want)
	}
}
