package chess

import "testing"

func mustPlay(t *testing.T, fen string, moves ...string) *Snapshot {
	t.Helper()
	g, err := newGame(fen)
	if err != nil {
		t.Fatalf("newGame(%q): %v", fen, err)
	}
	for _, m := range moves {
		if _, ok := playSAN(g, m); !ok {
			t.Fatalf("expected %s to be legal", m)
		}
	}
	s := snapshot(g)
	return &s
}

func TestStartPosition(t *testing.T) {
	g, err := newGame(StandardFEN)
	if err != nil {
		t.Fatal(err)
	}
	if got := currentFEN(g); got != StandardFEN {
		t.Errorf("fen = %q, want %q", got, StandardFEN)
	}
	s := snapshot(g)
	if s.MoveNumber != 0 || s.SAN != "" || s.From != nil || s.To != nil {
		t.Errorf("start snapshot should have no last move: %+v", s)
	}
	if s.SideToMove != White {
		t.Errorf("side to move = %v, want White", s.SideToMove)
	}
	if p := s.Board[4][0]; p == nil || *p != (Piece{White, King}) {
		t.Errorf("e1 = %v, want white king", p)
	}
}

func TestSnapshotReflectsLastMove(t *testing.T) {
	s := mustPlay(t, StandardFEN, "e4", "e5")
	if s.SAN != "e5" || s.MoveNumber != 2 || s.SideToMove != White {
		t.Errorf("unexpected snapshot: %+v", s)
	}
	if *s.From != (Square{4, 6}) || *s.To != (Square{4, 4}) {
		t.Errorf("from/to = %v/%v, want e7/e5", *s.From, *s.To)
	}
}

func TestPlaySANReturnsCanonicalNotation(t *testing.T) {
	g, _ := newGame(StandardFEN)
	for _, m := range []string{"f3", "e5", "g4"} {
		playSAN(g, m)
	}
	got, ok := playSAN(g, "Qh4")
	if !ok || got != "Qh4#" {
		t.Errorf("playSAN(Qh4) = %q, %v; want Qh4#, true", got, ok)
	}
}

func TestIllegalMovesAreRejected(t *testing.T) {
	g, _ := newGame(StandardFEN)
	for _, m := range []string{"e5", "Ke2", "zz", "", "  "} {
		if _, ok := playSAN(g, m); ok {
			t.Errorf("%q should be rejected", m)
		}
	}
	if got := currentFEN(g); got != StandardFEN {
		t.Errorf("position changed after rejected moves: %s", got)
	}
}

func TestMoveLeavingKingInCheckIsRejected(t *testing.T) {
	// White king e1 is in check from the rook on e8; the pawn move does not resolve it.
	g, _ := newGame("4r2k/8/8/8/8/8/P7/4K3 w - - 0 1")
	if _, ok := playSAN(g, "a3"); ok {
		t.Error("a3 should be rejected")
	}
	if _, ok := playSAN(g, "Kd1"); !ok {
		t.Error("Kd1 should be legal")
	}
}

func TestCheckAndCheckmateAreReported(t *testing.T) {
	s := mustPlay(t, StandardFEN, "e4", "f5", "Qh5")
	if !s.IsCheck || s.IsCheckmate {
		t.Errorf("Qh5: check=%v mate=%v, want check only", s.IsCheck, s.IsCheckmate)
	}
	s = mustPlay(t, StandardFEN, "f3", "e5", "g4", "Qh4")
	if !s.IsCheck || !s.IsCheckmate {
		t.Errorf("Qh4: check=%v mate=%v, want both", s.IsCheck, s.IsCheckmate)
	}
}

func TestCastling(t *testing.T) {
	s := mustPlay(t, "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", "O-O")
	if p := s.Board[6][0]; p == nil || *p != (Piece{White, King}) {
		t.Errorf("g1 = %v, want white king", p)
	}
	if p := s.Board[5][0]; p == nil || *p != (Piece{White, Rook}) {
		t.Errorf("f1 = %v, want white rook", p)
	}
}

func TestPromotionCanChooseAnyPiece(t *testing.T) {
	s := mustPlay(t, "7k/P7/8/8/8/8/8/K7 w - - 0 1", "a8=N")
	if p := s.Board[0][7]; p == nil || *p != (Piece{White, Knight}) {
		t.Errorf("a8 = %v, want white knight", p)
	}
}

func TestFENStartsAtThatPosition(t *testing.T) {
	const fen = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"
	g, err := newGame(fen)
	if err != nil {
		t.Fatal(err)
	}
	if got := currentFEN(g); got != fen {
		t.Errorf("fen = %q, want %q", got, fen)
	}
	if snapshot(g).SideToMove != Black {
		t.Error("side to move should be Black")
	}
}

func TestInvalidFENIsRejected(t *testing.T) {
	for _, fen := range []string{"", "bad", "8/8/8 w - - 0 1"} {
		if _, err := newGame(fen); err == nil {
			t.Errorf("newGame(%q) should fail", fen)
		}
	}
}

func TestFENWithWrongPiecesIsRejected(t *testing.T) {
	for _, fen := range []string{
		"8/8/8/8/8/8/8/8 w - - 0 1",
		"8/8/8/8/8/8/8/4K3 w - - 0 1",
		"4k3/8/8/8/8/8/8/8 w - - 0 1",
		"4k3/8/8/8/8/8/8/4K2K w - - 0 1",
		"k3k3/8/8/8/8/8/8/4K3 w - - 0 1",
		"P3k3/8/8/8/8/8/8/4K3 w - - 0 1",
		"4k3/8/8/8/8/8/8/p3K3 w - - 0 1",
	} {
		if _, err := newGame(fen); err == nil {
			t.Errorf("newGame(%q) should fail", fen)
		}
	}
}
