package chess

import (
	"slices"
	"testing"
)

func play(t *testing.T, moves ...string) *History {
	t.Helper()
	h, err := NewHistory(StandardFEN)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range moves {
		if !h.TryPlay(m) {
			t.Fatalf("expected %s to be legal", m)
		}
	}
	return h
}

func TestNewHistoryStartsAtStartPosition(t *testing.T) {
	h := play(t)

	if h.Position() != 0 || len(h.Moves()) != 0 || h.CanGoBack() || h.CanGoForward() {
		t.Errorf("unexpected initial state: pos=%d moves=%v", h.Position(), h.Moves())
	}
	if h.CurrentFEN() != StandardFEN {
		t.Errorf("fen = %q", h.CurrentFEN())
	}
}

func TestTryPlayAppendsMoveAndAdvances(t *testing.T) {
	h := play(t, "e4", "e5")

	if h.Position() != 2 || !slices.Equal(h.Moves(), []string{"e4", "e5"}) {
		t.Errorf("pos=%d moves=%v", h.Position(), h.Moves())
	}
}

func TestBackAndForwardMoveAlongTheLine(t *testing.T) {
	h := play(t, "e4", "e5", "Nf3")

	h.Back()
	if h.Position() != 2 || !h.CanGoForward() {
		t.Errorf("after Back: pos=%d forward=%v", h.Position(), h.CanGoForward())
	}
	h.Back()
	h.Back()
	if h.Position() != 0 || h.CanGoBack() || h.CurrentFEN() != StandardFEN {
		t.Errorf("at start: pos=%d fen=%s", h.Position(), h.CurrentFEN())
	}
	h.Forward()
	if h.Position() != 1 || h.Snapshot().SAN != "e4" {
		t.Errorf("after Forward: pos=%d san=%q", h.Position(), h.Snapshot().SAN)
	}
}

func TestBackAtStartAndForwardAtEndStayPut(t *testing.T) {
	h := play(t, "e4")

	h.Forward()
	if h.Position() != 1 {
		t.Errorf("Forward at end: pos=%d", h.Position())
	}
	h.Back()
	h.Back()
	if h.Position() != 0 {
		t.Errorf("Back at start: pos=%d", h.Position())
	}
}

func TestGoToClampsToValidRange(t *testing.T) {
	h := play(t, "e4", "e5")

	h.GoTo(99)
	if h.Position() != 2 {
		t.Errorf("GoTo(99): pos=%d", h.Position())
	}
	h.GoTo(-5)
	if h.Position() != 0 {
		t.Errorf("GoTo(-5): pos=%d", h.Position())
	}
}

func TestPlayingAfterGoingBackDiscardsMovesAhead(t *testing.T) {
	h := play(t, "e4", "e5", "Nf3")
	h.GoTo(1)

	if !h.TryPlay("c5") {
		t.Fatal("c5 should be legal")
	}

	if !slices.Equal(h.Moves(), []string{"e4", "c5"}) || h.Position() != 2 || h.CanGoForward() {
		t.Errorf("pos=%d moves=%v", h.Position(), h.Moves())
	}
}

func TestIllegalMoveIsRejectedAndNothingChanges(t *testing.T) {
	h := play(t, "e4")
	fen := h.CurrentFEN()

	if h.TryPlay("Ke5") {
		t.Error("Ke5 should be rejected")
	}
	if h.Position() != 1 || h.CurrentFEN() != fen || len(h.Moves()) != 1 {
		t.Errorf("state changed: pos=%d moves=%v", h.Position(), h.Moves())
	}
}

func TestIllegalMoveAfterGoingBackDoesNotDiscardMovesAhead(t *testing.T) {
	h := play(t, "e4", "e5", "Nf3")
	h.GoTo(1)

	if h.TryPlay("Ke5") {
		t.Error("Ke5 should be rejected")
	}
	if len(h.Moves()) != 3 || h.Position() != 1 || !h.CanGoForward() {
		t.Errorf("pos=%d moves=%v", h.Position(), h.Moves())
	}
}

func TestSnapshotAtStartHasNoLastMove(t *testing.T) {
	h := play(t, "e4")
	h.GoTo(0)

	s := h.Snapshot()
	if s.SAN != "" || s.From != nil {
		t.Errorf("unexpected snapshot: %+v", s)
	}
}

func TestCheckmateIsClearedWhenSteppingBack(t *testing.T) {
	h := play(t, "f3", "e5", "g4", "Qh4")
	if !h.Snapshot().IsCheckmate {
		t.Fatal("expected checkmate at the end of the line")
	}

	h.Back()
	if s := h.Snapshot(); s.IsCheckmate || s.IsCheck {
		t.Errorf("one move earlier: check=%v mate=%v", s.IsCheck, s.IsCheckmate)
	}
}

func TestNewHistoryFromFEN(t *testing.T) {
	const fen = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"

	h, err := NewHistory(fen)
	if err != nil {
		t.Fatal(err)
	}
	if h.StartFEN() != fen || h.Snapshot().SideToMove != Black {
		t.Errorf("start=%q side=%v", h.StartFEN(), h.Snapshot().SideToMove)
	}
}

func TestNewHistoryRejectsInvalidFEN(t *testing.T) {
	for _, fen := range []string{"", "bad", "8/8/8 w - - 0 1"} {
		if h, err := NewHistory(fen); err == nil || h != nil {
			t.Errorf("NewHistory(%q) should fail", fen)
		}
	}
}

func TestRestoreRebuildsLineAndPosition(t *testing.T) {
	h, err := RestoreHistory(StandardFEN, []string{"e4", "e5", "Nf3"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Moves()) != 3 || h.Position() != 2 {
		t.Errorf("pos=%d moves=%v", h.Position(), h.Moves())
	}
}

func TestRestoreRejectsIllegalMove(t *testing.T) {
	if _, err := RestoreHistory(StandardFEN, []string{"e4", "Ke5"}, 2); err == nil {
		t.Error("expected an error")
	}
}

func TestRestoreRejectsPositionOutOfRange(t *testing.T) {
	for _, position := range []int{-1, 3} {
		if _, err := RestoreHistory(StandardFEN, []string{"e4", "e5"}, position); err == nil {
			t.Errorf("position %d should be rejected", position)
		}
	}
}

func TestPreviewAgreesWithWhatEnterWouldPlay(t *testing.T) {
	for _, input := range []string{"e4", " Nf3", "Qh9", ""} {
		got := play(t).Preview(input).State == PreviewLegal
		want := play(t).TryPlay(input)

		if got != want {
			t.Errorf("%q: preview legal = %v, Enter plays = %v", input, got, want)
		}
	}
}

func TestPreviewReportsTheSquaresAndLeavesTheHistoryAlone(t *testing.T) {
	h := play(t)

	p := h.Preview("Nf3")

	if p.From == nil || *p.From != (Square{6, 0}) || p.To == nil || *p.To != (Square{5, 2}) {
		t.Errorf("from/to = %v/%v, want g1/f3", p.From, p.To)
	}
	if len(h.Moves()) != 0 || h.Position() != 0 {
		t.Errorf("history changed: moves=%v position=%d", h.Moves(), h.Position())
	}
}

func TestPreviewTellsLegalPartialAndInvalidInput(t *testing.T) {
	for input, want := range map[string]PreviewState{
		"Nf3": PreviewLegal,
		"N":   PreviewPartial,
		"Nf":  PreviewPartial,
		"e2":  PreviewPartial,
		"Nf5": PreviewInvalid,
		"":    PreviewNone,
	} {
		if got := play(t).Preview(input).State; got != want {
			t.Errorf("Preview(%q) = %v, want %v", input, got, want)
		}
	}
}
